package http

import (
	"context"
	"sync"
	"net/http"

	"github.com/bd878/gallery/server/logger"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	repository "github.com/bd878/gallery/server/users/internal/repository/sqlite"
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
	TableName           string
	DBPath              string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) (server *Server) {
	mux := http.NewServeMux()

	messagesGateway := messagesgateway.New(cfg.MessagesServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	repo := repository.New(cfg.TableName, cfg.DBPath)
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

	server = &Server{
		Server: &http.Server{
			Addr: cfg.Addr,
			Handler: mux,
		},
		config: cfg,
	}

	return
}

func (s *Server) ListenAndServe(_ context.Context, wg *sync.WaitGroup) {
	s.Server.ListenAndServe()
	wg.Done()
}