package http

import (
  "context"
  "sync"
  "net/http"

  httpmiddleware "github.com/bd878/gallery/server/messages/internal/middleware/http"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
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

  grpcCtrl := controller.New(controller.Config{
    RpcAddr: cfg.RpcAddr,
  })
  userGateway := usergateway.New(cfg.UsersServiceAddr)
  handler := httphandler.New(grpcCtrl, userGateway, cfg.DataPath)

  mux.Handle("/messages/v1/send", http.HandlerFunc(httpmiddleware.Logging(handler.CheckAuth(handler.SendMessage))))
  mux.Handle("/messages/v1/read", http.HandlerFunc(httpmiddleware.Logging(handler.CheckAuth(handler.ReadMessages))))
  mux.Handle("/messages/v1/status", http.HandlerFunc(httpmiddleware.Logging(handler.GetStatus)))
  mux.Handle("/messages/v1/read_file", http.HandlerFunc(httpmiddleware.Logging(handler.CheckAuth(handler.ReadFile))))

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