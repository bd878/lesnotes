package grpc

import (
	"net"
	"time"
	"os"
	"fmt"
	"sync"
	"context"
	"google.golang.org/grpc"
	"github.com/soheilhy/cmux"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"

	"github.com/bd878/gallery/server/waiter"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	controller "github.com/bd878/gallery/server/sessions/internal/controller/sessions"
	repository "github.com/bd878/gallery/server/sessions/internal/repository/postgres"
	grpchandler "github.com/bd878/gallery/server/sessions/internal/handler/grpc"
)

type Config struct {
	Addr            string
	PGConn          string
	NodeName        string
}

type Server struct {
	*grpc.Server
	conf             Config
	mux              cmux.CMux
	pool             *pgxpool.Pool
	listener         net.Listener
	grpcListener     net.Listener
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
	server.setupGRPC(logger.Default())

	return server
}

func (s *Server) setupGRPC(log *logger.Logger) {
	repo := repository.New(s.pool)
	control := controller.New(repo)
	handler := grpchandler.New(control)

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
	)

	api.RegisterSessionsServer(s.Server, handler)

	s.grpcListener = s.mux.Match(cmux.Any())
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForRPC, s.WaitForPool)

	return waiter.Wait()
}

func (s *Server) setupDB() {
	var err error
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	if err != nil {
		panic(err)
	}
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