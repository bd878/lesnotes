package grpc

import (
	"net"
	"os"
	"io"
	"bytes"
	"fmt"
	"time"
	"context"
	"google.golang.org/grpc"
	"github.com/soheilhy/cmux"
	"github.com/hashicorp/raft"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/bd878/gallery/server/api"

	"github.com/bd878/gallery/server/internal/waiter"
	membership "github.com/bd878/gallery/server/internal/discovery/serf"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/distributed"
	repository "github.com/bd878/gallery/server/users/internal/repository/postgres"
	grpchandler "github.com/bd878/gallery/server/users/internal/handler/grpc"
)

type Config struct {
	Addr                   string
	PGConn                 string
	TableName              string
	RaftLogLevel           string
	RaftBootstrap          bool
	RaftServers            []string
	SerfAddr               string
	SerfJoinAddrs          []string
	NodeName               string
	DataPath               string
}

type Server struct {
	*grpc.Server
	conf             Config
	mux              cmux.CMux
	pool             *pgxpool.Pool
	listener         net.Listener
	grpcListener     net.Listener
	controller       *controller.DistributedUsers
	membership       *membership.Membership
}

func New(cfg Config) *Server {
	listener, err := net.Listen("tcp4", cfg.Addr)
	if err != nil {
		panic(err)
	}

	mux := cmux.New(listener)

	server := &Server{
		conf:          cfg,
		mux:           mux,
		listener:      listener,
	}

	if err := server.setupDB(); err != nil {
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

func (s *Server) setupGRPC() (err error) {
	s.membership, err = membership.New(
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
	go s.membership.Run(context.Background())

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
	)

	api.RegisterUsersServer(s.Server, grpchandler.New(s.controller))

	s.grpcListener = s.mux.Match(cmux.Any())

	return nil
}

func (s *Server) setupDB() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	return
}

func (s *Server) setupRaft() (err error) {
	repo := repository.New(s.pool, s.conf.TableName)

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
			LocalID:  raft.ServerID(s.conf.NodeName),
			LogLevel: raftLogLevel,
		},
		StreamLayer: controller.NewStreamLayer(raftListener),
		Bootstrap:   s.conf.RaftBootstrap,
		DataDir:     s.conf.DataPath,
		Servers:     s.conf.RaftServers,
	}, repo)

	return
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForRPC, s.WaitForPool)

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

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}