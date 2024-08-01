package main

import (
  "fmt"
  "net/http"

  config "github.com/bd878/gallery/server/messages/config"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  controller "github.com/bd878/gallery/server/messages/internal/controller/service"
)

func New(cfg config.Config) *http.Server {
  mux := http.NewServeMux()

  ctrlCfg := controller.Config{
    ClusterNodeAddr: cfg.LeaderAddr,
  }

  grpcCtrl := controller.New(ctrlCfg)
  userGateway := usergateway.New(cfg.UserAddr)
  h := httphandler.New(grpcCtrl, userGateway, cfg.DataPath)

  mux.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SendMessage)))
  mux.Handle("/messages/v1/read", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  mux.Handle("/messages/v1/status", http.HandlerFunc(h.GetStatus))
  mux.Handle("/messages/v1/read_file", http.HandlerFunc(h.CheckAuth(h.ReadFile)))

  srv := &http.Server{
    Addr: fmt.Sprintf(":%d", cfg.Port),
    Handler: mux,
  }

  return srv
}