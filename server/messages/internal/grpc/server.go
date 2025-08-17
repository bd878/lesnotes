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

	"golang.org/x/sync/errgroup"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	hclog "github.com/hashicorp/go-hclog"

	"github.com/bd878/gallery/server/waiter"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	grpclogger "github.com/bd878/gallery/server/messages/internal/logger/grpc"
	membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
	repository "github.com/bd878/gallery/server/messages/internal/repository/postgres"
	controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
	grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type Config struct {
	Addr                string
	PGConn              string
	DBPath              string
	NodeName            string
	RaftLogLevel        string
	RaftBootstrap       bool
	DataPath            string
	RaftServers         []string
	SerfAddr            string
	SerfJoinAddrs       []string
	SessionsServiceAddr string
}

type Server struct {
	*grpc.Server
	conf             Config
	mux              cmux.CMux
	pool             *pgxpool.Pool
	listener         net.Listener
	grpcListener     net.Listener
	controller      *controller.DistributedMessages
	membership      *membership.Membership
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

	server.setupDB()
	server.setupRaft(logger.Default())
	server.setupGRPC(logger.Default())

	return server
}

func (s *Server) setupDB() {
	var err error
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	if err != nil {
		panic(err)
	}
}

func (s *Server) setupRaft(log *logger.Logger) {
	repo := repository.New(s.pool)

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

	control, err := controller.New(repo, controller.Config{
		Raft: raft.Config{
			LocalID: raft.ServerID(s.conf.NodeName),
			LogLevel: raftLogLevel,
		},
		StreamLayer: controller.NewStreamLayer(raftListener),
		Bootstrap:   s.conf.RaftBootstrap,
		DataDir:     s.conf.DataPath,
		Servers:     s.conf.RaftServers,
	})
	if err != nil {
		panic(err)
	}

	s.controller = control
}

func (s *Server) setupGRPC(log *logger.Logger) {
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
		log.Errorw("failed to establish membership connection", "error", err)
		panic(err)
	}

	s.membership = member

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpclogger.NewBuilder()),
		),
	)

	api.RegisterMessagesServer(s.Server, handler)

	grpcListener := s.mux.Match(cmux.Any())
	s.grpcListener = grpcListener
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForRPC, s.WaitForPool)

	return waiter.Wait()
}

func (s *Server) WaitForPool(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", s.Addr())
		defer fmt.Fprintf(os.Stdout, "rpc server shutdown")
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
		s.mux.Serve()
		defer s.mux.Close()
		return nil
	})

	return group.Wait()
}

func (s *Server) WaitForRPC(ctx context.Context) (err error) {
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
