package http

import (
	"context"
	"sync"
	"net/http"

	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
	httplogger "github.com/bd878/gallery/server/messages/internal/logger/http"
	usergateway "github.com/bd878/gallery/server/internal/gateway/user"
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

	userGateway := usergateway.New(cfg.UsersServiceAddr)
	filesGateway := filesgateway.New(cfg.FilesServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(userGateway, usermodel.PublicUserID))

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

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(userGateway, usermodel.PublicUserID))
	mux.Handle("/messages/v2/send", middleware.Build(handler.SendMessageJsonAPI))

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