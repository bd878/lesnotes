package grpc

import (
	"net"
	"fmt"
	"os"
	"time"
	"context"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"github.com/nats-io/nats.go"

	"github.com/bd878/gallery/server/api"

	"github.com/bd878/gallery/server/waiter"
	"github.com/bd878/gallery/server/ddd"
	broker "github.com/bd878/gallery/server/nats"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	repository "github.com/bd878/gallery/server/files/internal/repository/postgres"
	streamhandler "github.com/bd878/gallery/server/files/internal/handler/stream"
	grpchandler "github.com/bd878/gallery/server/files/internal/handler/grpc"
)

type Config struct {
	Addr                  string
	PGConn                string
	NodeName              string
	SessionsServiceAddr   string
	NatsAddr              string
}

type Server struct {
	*grpc.Server
	conf             Config
	nc               *nats.Conn
	pool             *pgxpool.Pool
	listener         net.Listener
}

func New(cfg Config) *Server {
	listener, err := net.Listen("tcp4", cfg.Addr)
	if err != nil {
		panic(err)
	}

	server := &Server{
		conf:          cfg,
		listener:      listener,
	}

	if err := server.setupDB(); err != nil {
		panic(err)
	}

	if err := server.setupNats(); err != nil {
		panic(err)
	}

	if err := server.setupGRPC(); err != nil {
		panic(err)
	}

	return server
}

func (s *Server) setupGRPC() (err error) {
	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream := broker.NewStream(s.nc)
	streamhandler.RegisterDomainEventHandlers(dispatcher, streamhandler.NewDomainEventHandlers(stream))

	repo := repository.New(s.pool)
	handler := grpchandler.New(repo, dispatcher)

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
		grpc.MaxRecvMsgSize(1024*1024*50),
		grpc.MaxSendMsgSize(1024*1024*50),
	)

	api.RegisterFilesServer(s.Server, handler)

	return nil
}

func (s *Server) setupNats() (err error) {
	s.nc, err = nats.Connect(s.conf.NatsAddr)

	return
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForRPC, s.WaitForPool, s.WaitForStream)

	return waiter.Wait()
}

func (s *Server) setupDB() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	return
}

func (s *Server) WaitForRPC(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", s.Addr())
		defer fmt.Fprintln(os.Stdout, "rpc server shutdown")
		if err := s.Serve(s.listener); err != nil {
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
		fmt.Fprintln(os.Stdout, "message stream started")
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