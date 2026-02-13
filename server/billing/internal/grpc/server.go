package grpc

import (
	"os"
	"net"
	"io"
	"fmt"
	"time"
	"context"
	"bytes"
	"google.golang.org/grpc"
	"github.com/hashicorp/raft"
	"github.com/soheilhy/cmux"
	"github.com/nats-io/nats.go"

	"golang.org/x/sync/errgroup"
	"github.com/bd878/gallery/server/api"
	"github.com/jackc/pgx/v5/pgxpool"
	hclog "github.com/hashicorp/go-hclog"

	"github.com/bd878/gallery/server/waiter"
	membership "github.com/bd878/gallery/server/discovery/serf"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	repository "github.com/bd878/gallery/server/billing/internal/repository/postgres"
	controller "github.com/bd878/gallery/server/billing/internal/controller/distributed"
	grpchandler "github.com/bd878/gallery/server/billing/internal/handler/grpc"
)

type Config struct {
	Addr                string
	PGConn              string
	PaymentsTableName   string
	InvoicesTableName   string
	NodeName            string
	RaftLogLevel        string
	RaftBootstrap       bool
	DataPath            string
	RaftServers         []string
	SerfAddr            string
	SerfJoinAddrs       []string
	NatsAddr            string
}

type Server struct {
	*grpc.Server
	conf             Config
	nc               *nats.Conn
	mux              cmux.CMux
	pool             *pgxpool.Pool
	listener         net.Listener
	grpcListener     net.Listener
	controller       *controller.Distributed
	membership       *membership.Membership
}

func New(cfg Config) *Server {
	listener, err := net.Listen("tcp4", cfg.Addr)
	if err != nil {
		panic(err)
	}

	mux := cmux.New(listener)

	server := &Server{
		conf:        cfg,
		mux:         mux,
		listener:    listener,
	}

	if err := server.setupDB(); err != nil {
		panic(err)
	}

	if err := server.setupNats(); err != nil {
		panic(err)
	}

	if err := server.setupRaft(); err != nil {
		panic(err)
	}

	if err := server.setupGRPC(); err != nil {
		panic(err)
	}

	return server
}

func (s *Server) setupDB() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	return
}

func (s *Server) setupNats() (err error) {
	s.nc, err = nats.Connect(s.conf.NatsAddr)
	return
}

func (s *Server) setupRaft() (err error) {
	paymentsRepo := repository.NewPaymentsRepository(s.pool, s.conf.PaymentsTableName)
	invoicesRepo := repository.NewInvoicesRepository(s.pool, s.conf.InvoicesTableName)

	raftLogLevel := hclog.Error.String()
	switch s.conf.RaftLogLevel {
	case "debug":
		raftLogLevel = hclog.Debug.String()
	case "error":
		raftLogLevel = hclog.Error.String()
	case "info":
		raftLogLevel = hclog.Info.String()
	default:
		raftLogLevel = hclog.Info.String()
	}

	raftListener := s.mux.Match(func(r io.Reader) bool {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(controller.RaftRPC)}) == 0
	})

	s.controller, err = controller.New(controller.Config{
		Raft: raft.Config{
			LocalID: raft.ServerID(s.conf.NodeName),
			LogLevel: raftLogLevel,
		},
		StreamLayer: controller.NewStreamLayer(raftListener),
		Bootstrap:   s.conf.RaftBootstrap,
		DataDir:     s.conf.DataPath,
		Servers:     s.conf.RaftServers,
	}, paymentsRepo, invoicesRepo)

	return
}

func (s *Server) setupGRPC() error {
	handler := grpchandler.New(s.controller)
	member, err := membership.New(
		membership.Config{
			NodeName: s.conf.NodeName,
			BindAddr: s.conf.SerfAddr,
			Tags: map[string]string{
				"raft_addr": s.conf.Addr,
			},
			SerfJoinAddrs: s.conf.SerfJoinAddrs,
		},
		s.controller,
	)
	if err != nil {
		return err
	}

	s.membership = member

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
	)

	api.RegisterBillingServer(s.Server, handler)

	s.grpcListener = s.mux.Match(cmux.Any())

	return nil
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForRPC, s.WaitForPool, s.WaitForStream)

	return waiter.Wait()
}

func (s *Server) WaitForRPC(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", s.Addr())
		defer fmt.Fprintln(os.Stdout, "rpc server shutdown")
		if err := s.Serve(s.grpcListener); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(5*time.Second)
		select {
		case <-timeout.C:
			s.Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})
	group.Go(func() error {
		fmt.Fprintln(os.Stdout, "mux serve")
		s.mux.Serve()
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "mux close")
		s.mux.Close()
		return nil
	})

	return group.Wait()
}

func (s *Server) WaitForPool(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "closing pgpool connections")
		s.pool.Close()
		return nil
	})

	return group.Wait()
}

func (s *Server) WaitForStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func (*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Fprintln(os.Stdout, "messsage stream started")
		defer fmt.Fprintln(os.Stdout, "message stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})
	return group.Wait()
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}
