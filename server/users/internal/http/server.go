package http

import (
  "context"
  "sync"
  "net/http"

  httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
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

  middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

  repo := repository.New(cfg.DBPath)
  grpcCtrl := controller.New(repo)
  handler := httphandler.New(grpcCtrl, httphandler.Config{
    Domain:    cfg.Domain,
  })

  mux.Handle("/users/v1/signup", middleware.Build(handler.Register))
  mux.Handle("/users/v1/login",  middleware.Build(handler.Authenticate))
  mux.Handle("/users/v1/auth",   middleware.Build(handler.Auth))
  mux.Handle("/users/v1/status", middleware.Build(handler.ReportStatus))

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