package grpc

import (
  "net"
  "sync"
  "context"
  "google.golang.org/grpc"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/logger"

  grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
  repository "github.com/bd878/gallery/server/files/internal/repository/sqlite"
  grpchandler "github.com/bd878/gallery/server/files/internal/handler/grpc"
)

type Config struct {
  Addr            string
  DBPath          string
  NodeName        string
  DataPath        string
}

type Server struct {
  *grpc.Server
  config           Config
  listener         net.Listener
}

func New(cfg Config) *Server {
  listener, err := net.Listen("tcp4", cfg.Addr)
  if err != nil {
    panic(err)
  }

  server := &Server{
    config:        cfg,
    listener:      listener,
  }

  server.setupGRPC(logger.Default())

  return server
}

func (s *Server) setupGRPC(log *logger.Logger) {
  repo := repository.New(s.config.DBPath)
  handler := grpchandler.New(repo, s.config.DataPath)

  s.Server = grpc.NewServer(
    grpc.ChainUnaryInterceptor(
      grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
    ),
  )

  api.RegisterFilesServer(s.Server, handler)
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
  logger.Info("server is listening on ", s.config.Addr)
  s.Serve(s.listener)
  wg.Done()
}

func (s *Server) Addr() string {
  return s.listener.Addr().String()
}