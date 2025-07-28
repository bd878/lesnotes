package http

import (
	"context"
	"sync"
	"net/http"

	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/files/internal/handler/http"
	controller "github.com/bd878/gallery/server/files/internal/controller/service"
)

type Config struct {
	Addr             string
	RpcAddr          string
	UsersServiceAddr string
	DataPath         string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))

	grpcCtrl := controller.New(controller.Config{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(grpcCtrl)

	mux.Handle("/files/v1/upload", middleware.Build(handler.UploadFile))
	mux.Handle("/files/v1/download", middleware.Build(handler.DownloadFile))

	middleware.NoAuth()
	mux.Handle("/files/v1/status", middleware.Build(handler.GetStatus))
	mux.Handle("/files/v2/{user_id}/{name}", middleware.Build(handler.DownloadFileV2))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/files/v2/upload", middleware.Build(handler.UploadFileV2))

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
	err := s.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	wg.Done()
}