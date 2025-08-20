package http

import (
	"context"
	"fmt"
	"os"
	"time"
	"net/http"

	"golang.org/x/sync/errgroup"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/waiter"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	repository "github.com/bd878/gallery/server/users/internal/repository/postgres"
	httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
	messagesgateway "github.com/bd878/gallery/server/users/internal/gateway/messages/grpc"
	sessionsgateway "github.com/bd878/gallery/server/users/internal/gateway/sessions/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/users"
)

type Config struct {
	Addr                string
	RpcAddr             string
	MessagesServiceAddr string
	SessionsServiceAddr string
	DataPath            string
	CookieDomain        string
	PGConn              string
}

type Server struct {
	*http.Server
	conf             Config
	pool             *pgxpool.Pool
}

func New(cfg Config) (server *Server) {
	mux := http.NewServeMux()

	server = &Server{
		Server: &http.Server{
			Addr: cfg.Addr,
			Handler: mux,
		},
		conf: cfg,
	}

	server.setupDB()

	messagesGateway := messagesgateway.New(cfg.MessagesServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	repo := repository.New(server.pool)
	control := controller.New(repo, messagesGateway, sessionsGateway)
	handler := httphandler.New(control, httphandler.Config{
		CookieDomain:    cfg.CookieDomain,
	})

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), control, sessionsGateway, usersmodel.PublicUserID))
	mux.Handle("/users/v1/get",    middleware.Build(handler.GetUser))
	mux.Handle("/users/v1/logout", middleware.Build(handler.Logout))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), control, sessionsGateway, usersmodel.PublicUserID))
	mux.Handle("/users/v2/delete", middleware.Build(handler.DeleteJsonAPI))

	middleware.NoAuth()
	mux.Handle("/users/v1/signup", middleware.Build(handler.Signup))
	mux.Handle("/users/v1/login",  middleware.Build(handler.Login))
	mux.Handle("/users/v1/auth",   middleware.Build(handler.Auth))
	mux.Handle("/users/v1/status", middleware.Build(handler.Status))
	mux.Handle("/users/v2/signup", middleware.Build(handler.SignupJsonAPI))
	mux.Handle("/users/v2/auth",   middleware.Build(handler.AuthJsonAPI))
	mux.Handle("/users/v2/login",  middleware.Build(handler.LoginJsonAPI))

	return
}

func (s *Server) setupDB() {
	var err error
	s.pool, err = pgxpool.New(context.Background(), s.conf.PGConn)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForServer, s.WaitForPool)

	return waiter.Wait()
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


func (s *Server) WaitForServer(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "http server started %s\n", s.Addr)
		defer fmt.Fprintf(os.Stdout, "http server shutdown")
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
