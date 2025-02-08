package http

import (
  "context"
  "sync"
  "net/http"

  httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
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

  middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

  userGateway := usergateway.New(cfg.UsersServiceAddr)
  filesGateway := filesgateway.New(cfg.FilesServiceAddr)
  authBuilder := &httpmiddleware.AuthBuilder{Gateway: userGateway}
  middleware = middleware.WithAuth(authBuilder.Auth)

  grpcCtrl := controller.New(controller.Config{
    RpcAddr: cfg.RpcAddr,
  })
  handler := httphandler.New(grpcCtrl, filesGateway)

  mux.Handle("/messages/v1/send", middleware.Build(handler.SendMessage))
  mux.Handle("/messages/v1/read", middleware.Build(handler.ReadMessages))
  mux.Handle("/messages/v1/update", middleware.Build(handler.UpdateMessage))
  mux.Handle("/messages/v1/delete", middleware.Build(handler.DeleteMessage))
  mux.Handle("/messages/v1/status", middleware.NoAuth().Build(handler.GetStatus))

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