package http

import (
  "context"
  "sync"
  "net/http"

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
  authBuilder := &httpmiddleware.AuthBuilder{Gateway: userGateway}
  middleware = middleware.WithAuth(authBuilder.Auth)

  grpcCtrl := controller.New(controller.Config{
    RpcAddr: cfg.RpcAddr,
  })
  handler := httphandler.New(grpcCtrl)

  mux.Handle("/files/v1/{file_id}", middleware.Build(handler.DownloadFile))
  mux.Handle("/files/v1/status", middleware.NoAuth().Build(handler.GetStatus))

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