package main

import (
  "context"
  "sync"
  "net/http"

  config "github.com/bd878/gallery/server/messages/config"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  controller "github.com/bd878/gallery/server/messages/internal/controller/service"
)

type HTTPServer struct {
  *http.Server
  config HTTPServerConfig
}

type HTTPServerConfig struct {
  Addr            string
  RpcAddr         string
  UserServiceAddr string
  DataPath        string
}

func NewHTTPServer(cfg HTTPServerConfig) *HTTPServer {
  mux := http.NewServeMux()

  grpcCtrl := controller.New(controller.Config{
    RpcAddr: cfg.RpcAddr,
  })
  userGateway := usergateway.New(cfg.UsersServiceAddr)
  h := httphandler.New(grpcCtrl, userGateway, cfg.DataPath)

  mux.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SendMessage)))
  mux.Handle("/messages/v1/read", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  mux.Handle("/messages/v1/status", http.HandlerFunc(h.GetStatus))
  mux.Handle("/messages/v1/read_file", http.HandlerFunc(h.CheckAuth(h.ReadFile)))

  server := &HTTPServer{
    Server: &http.Server{
      Addr: cfg.Addr,
      Handler: mux,
    },
    config: cfg,
  }

  return server
}

func (s *HTTPServer) ListenAndServe(_ context.Context, wg *sync.WaitGroup) {
  s.Server.ListenAndServe()
  wg.Done()
}