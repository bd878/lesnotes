package http

import (
	"context"
	"fmt"
	"os"
	"time"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/waiter"
	users "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
	sessionsgateway "github.com/bd878/gallery/server/users/internal/gateway/sessions/grpc"
	messagesgateway "github.com/bd878/gallery/server/users/internal/gateway/messages/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/service"
)

type Config struct {
	Addr                string
	RpcAddr             string
	SessionsServiceAddr string
	MessagesServiceAddr string
	CookieDomain        string
}

type Server struct {
	*http.Server
	conf  Config
}

func New(cfg Config) (server *Server) {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log).WithLang(httpmiddleware.Language)

	server = &Server{
		Server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
		conf: cfg,
	}

	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	messagesGateway := messagesgateway.New(cfg.MessagesServiceAddr)

	ctrl := controller.New(controller.Config{RpcAddr: cfg.RpcAddr}, messagesGateway, sessionsGateway)

	handler := httphandler.New(ctrl, httphandler.Config{
		CookieDomain:    cfg.CookieDomain,
	})

	// TODO: middleware.Build(handler, ...middlewares)
	middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), ctrl, sessionsGateway, users.PublicUserID))
	mux.Handle("/users/v1/me",     middleware.Build(handler.GetMe))
	mux.Handle("/users/v1/logout", middleware.Build(handler.Logout))
	mux.Handle("/users/v1/update", middleware.Build(handler.Update))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), ctrl, sessionsGateway, users.PublicUserID))
	mux.Handle("/users/v2/delete", middleware.Build(handler.DeleteJsonAPI))
	mux.Handle("/users/v2/me",     middleware.Build(handler.GetMe))
	mux.Handle("/users/v2/update", middleware.Build(handler.UpdateJsonAPI))

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

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForServer)

	return waiter.Wait()
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
