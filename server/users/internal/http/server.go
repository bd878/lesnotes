package http

import (
  "context"
  "sync"
  "net/http"

  repository "github.com/bd878/gallery/server/users/internal/repository/sqlite"
  httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
  controller "github.com/bd878/gallery/server/users/internal/controller/users"
)

type Config struct {
  Addr             string
  RpcAddr          string
  DataPath         string
  Domain           string
  DBPath           string
}

type Server struct {
  *http.Server
  config Config
}

func New(cfg Config) *Server {
  mux := http.NewServeMux()

  repo := repository.New(cfg.DBPath)
  grpcCtrl := controller.New(repo)
  handler := httphandler.New(grpcCtrl, httphandler.Config{
    Domain:    cfg.Domain,
  })

  mux.Handle("/users/v1/signup", http.HandlerFunc(handler.Register))
  mux.Handle("/users/v1/login",  http.HandlerFunc(handler.Authenticate))
  mux.Handle("/users/v1/auth",   http.HandlerFunc(handler.Auth))
  mux.Handle("/users/v1/status", http.HandlerFunc(handler.ReportStatus))

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