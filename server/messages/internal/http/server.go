package http

import (
  "context"
  "sync"
  "net/http"

  httpmiddleware "github.com/bd878/gallery/server/messages/internal/middleware/http"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/mock"
  controller "github.com/bd878/gallery/server/messages/internal/controller/service"
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
  handler := httphandler.New(grpcCtrl, cfg.DataPath)

  mux.Handle("/messages/v1/send", middleware.Build(handler.SendMessage))
  mux.Handle("/messages/v1/read", middleware.Build(handler.ReadMessages))
  mux.Handle("/messages/v1/update", middleware.Build(handler.UpdateMessage))
  mux.Handle("/messages/v1/status", middleware.NoAuth().Build(handler.GetStatus))
  mux.Handle("/messages/v1/read_file", middleware.Build(handler.ReadFile))

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