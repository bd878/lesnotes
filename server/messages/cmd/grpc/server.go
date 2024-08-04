package main

import (
  "net"
  "google.golang.org/grpc"
  "github.com/hashicorp/raft"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/config"

  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type GRPCMessagesServer struct {
  ln   net.Listener
  srv *grpc.Server
}

func New(cfg config.Config) *GRPCMessagesServer {
  ln, err := net.Listen("tcp4", cfg.GrpcAddr)
  if err != nil {
    panic(err)
  }

  streamLn, err := net.Listen("tcp", cfg.RaftAddr)
  if err != nil {
    panic(err)
  }

  repo, err := repository.New(cfg.DBPath)
  if err != nil {
    panic(err)
  }

  ctrl, err := controller.New(repo, controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(cfg.RaftAddr),
    },
    StreamLayer: controller.NewStreamLayer(streamLn),
    Bootstrap:   cfg.Bootstrap,
    DataDir:     cfg.DataPath,
    JoinAddrs:   cfg.JoinAddrs,
  })
  if err != nil {
    panic(err)
  }

  h := grpchandler.New(ctrl)

  srv := grpc.NewServer()
  api.RegisterMessagesServer(srv, h)

  a := &GRPCMessagesServer{
    ln: ln,
    srv: srv,
  }

  return a
}

func (s *GRPCMessagesServer) Run() {
  defer s.ln.Close()
  s.srv.Serve(s.ln)
}

func (s *GRPCMessagesServer) Addr() string {
  return s.ln.Addr().String()
}