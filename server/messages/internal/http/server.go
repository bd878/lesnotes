package http

import (
	"context"
	"sync"
	"net/http"

	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
	httplogger "github.com/bd878/gallery/server/messages/internal/logger/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	filesgateway "github.com/bd878/gallery/server/messages/internal/gateway/files/grpc"
	controller "github.com/bd878/gallery/server/messages/internal/controller/service"
)

type Config struct {
	Addr             string
	RpcAddr          string
	UsersServiceAddr string
	FilesServiceAddr string
	DataPath         string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httplogger.LogBuilder())

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	filesGateway := filesgateway.New(cfg.FilesServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))

	grpcCtrl := controller.New(controller.Config{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(grpcCtrl, filesGateway)

	mux.Handle("/messages/v1/send", middleware.Build(handler.SendMessage))
	mux.Handle("/messages/v1/read", middleware.Build(handler.ReadMessageOrMessages))
	mux.Handle("/messages/v1/update", middleware.Build(handler.UpdateMessage))
	mux.Handle("/messages/v1/publish", middleware.Build(handler.PublishMessageOrMessages))
	mux.Handle("/messages/v1/private", middleware.Build(handler.PrivateMessageOrMessages))
	mux.Handle("/messages/v1/delete", middleware.Build(handler.DeleteMessageOrMessages))

	middleware.NoAuth()
	mux.Handle("/messages/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/messages/v2/send", middleware.Build(handler.SendMessageJsonAPI))
	mux.Handle("/messages/v2/read", middleware.Build(handler.ReadMessageOrMessagesJsonAPI))
	mux.Handle("/messages/v2/publish", middleware.Build(handler.PublishMessageOrMessagesJsonAPI))
	mux.Handle("/messages/v2/private", middleware.Build(handler.PrivateMessageOrMessagesJsonAPI))
	mux.Handle("/messages/v2/delete", middleware.Build(handler.DeleteMessageOrMessagesJsonAPI))
	mux.Handle("/messages/v2/update", middleware.Build(handler.UpdateMessageJsonAPI))

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