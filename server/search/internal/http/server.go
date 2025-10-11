package http

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"golang.org/x/sync/errgroup"
	"github.com/nats-io/nats.go"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/waiter"
	"github.com/bd878/gallery/server/logger"
	broker "github.com/bd878/gallery/server/nats"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/search/internal/handler/http"
	streamhandler "github.com/bd878/gallery/server/search/internal/handler/stream"
	repository "github.com/bd878/gallery/server/search/internal/repository/postgres"
	controller "github.com/bd878/gallery/server/search/internal/controller/search"
)

type Config struct {
	Addr                string
	UsersServiceAddr    string
	SessionsServiceAddr string
	NatsAddr            string
	PGConn              string
	MessagesTableName   string
	FilesTableName      string
}

type Server struct {
	*http.Server
	conf   Config

	pool   *pgxpool.Pool
	nc     *nats.Conn
}

func New(conf Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	server := &Server{
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: mux,
		},
		conf: conf,
	}

	if err := server.setupNats(); err != nil {
		panic(err)
	}

	if err := server.setupDB(); err != nil {
		panic(err)
	}

	usersGateway := usersgateway.New(conf.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(conf.SessionsServiceAddr)

	repo := repository.New(server.pool, conf.MessagesTableName, conf.FilesTableName)
	ctrl := controller.New(controller.Config{}, repo, repo)

	streamhandler.RegisterIntegrationEventHandlers(broker.NewStream(server.nc), streamhandler.NewIntegrationEventHandlers(ctrl))

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/search/v1/messages", middleware.Build(handler.SearchMessages))
	mux.Handle("/search/v1/files", middleware.Build(handler.SearchFiles))

	middleware.NoAuth()
	mux.Handle("/search/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/search/v2/messages", middleware.Build(handler.SearchMessagesJsonAPI))
	mux.Handle("/search/v2/files", middleware.Build(handler.SearchFilesJsonAPI))

	return server
}

func (s *Server) setupNats() (err error) {
	s.nc, err = nats.Connect(s.conf.NatsAddr)

	return
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForPool, s.WaitForServer, s.WaitForStream)

	return waiter.Wait()
}

func (s *Server) setupDB() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	return
}

func (s *Server) WaitForServer(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "http server started %s\n", s.Addr)
		defer fmt.Fprintln(os.Stdout, "http server shutdown")
		if err := s.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "http server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "http server failed to stop gracefully")
			return err
		}
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
