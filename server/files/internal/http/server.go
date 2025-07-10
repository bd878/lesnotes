package http

import (
	"context"
	"sync"
	"net/http"

	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	usergateway "github.com/bd878/gallery/server/internal/gateway/user"
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

	userGateway := usergateway.New(cfg.UsersServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(userGateway, usermodel.PublicUserID))

	grpcCtrl := controller.New(controller.Config{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(grpcCtrl)

	mux.Handle("/files/v1/upload", middleware.Build(handler.UploadFile))
	mux.Handle("/files/v1/download", middleware.Build(handler.DownloadFile))

	middleware.NoAuth()
	mux.Handle("/files/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(userGateway, usermodel.PublicUserID))
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