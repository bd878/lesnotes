package http

import (
	"context"
	"sync"
	"net/http"

	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	repository "github.com/bd878/gallery/server/users/internal/repository/sqlite"
	httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
	messagesgateway "github.com/bd878/gallery/server/users/internal/gateway/messages/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/users"
)

type Config struct {
	Addr             string
	RpcAddr          string
	MessagesServiceAddr string
	DataPath         string
	CookieDomain     string
	DBPath           string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	messagesGateway := messagesgateway.New(cfg.MessagesServiceAddr)

	repo := repository.New(cfg.DBPath)
	grpcCtrl := controller.New(repo)
	handler := httphandler.New(grpcCtrl, httphandler.Config{
		CookieDomain:    cfg.CookieDomain,
	}, messagesGateway)

	mux.Handle("/users/v1/get", middleware.Build(handler.GetUser))
	mux.Handle("/users/v1/signup", middleware.Build(handler.Signup))
	mux.Handle("/users/v1/login",  middleware.Build(handler.Login))
	mux.Handle("/users/v1/logout", middleware.Build(handler.Logout))
	mux.Handle("/users/v1/auth",   middleware.Build(handler.Auth))
	mux.Handle("/users/v1/status", middleware.Build(handler.Status))
	mux.Handle("/users/v2/auth",   middleware.Build(handler.AuthJsonAPI))
	mux.Handle("/users/v2/login",  middleware.Build(handler.LoginJsonAPI))
	mux.Handle("/users/v2/signup", middleware.Build(handler.SignupJsonAPI))
	mux.Handle("/users/v2/delete", middleware.Build(handler.DeleteJsonAPI))

	server := &Server{
		Server: &http.Server{
			Addr: cfg.Addr,
			Handler: mux,
		},
		config: cfg,
	}

	return server
}

func (s *Server) ListenAndServe(_ context.Context, wg *sync.WaitGroup) {
	s.Server.ListenAndServe()
	wg.Done()
}