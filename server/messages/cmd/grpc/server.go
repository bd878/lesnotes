package main

import (
  "net"
  "google.golang.org/grpc"
  "github.com/hashicorp/raft"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/config"
  hclog "github.com/hashicorp/go-hclog"

  membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type GRPCMessagesServer struct {
  ln   net.Listener
  srv *grpc.Server
  m   *membership.Membership
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

  raftLogLevel := hclog.Error.String()
  switch cfg.RaftLogLevel {
  case "debug":
    raftLogLevel = hclog.Debug.String()
  case "error":
    raftLogLevel = hclog.Error.String()
  case "info":
    raftLogLevel = hclog.Info.String()
  }

  ctrl, err := controller.New(repo, controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(cfg.NodeName),
      LogLevel: raftLogLevel,
    },
    StreamLayer: controller.NewStreamLayer(streamLn),
    Bootstrap:   cfg.RaftBootstrap,
    DataDir:     cfg.DataPath,
    Servers:     cfg.RaftServers,
  })
  if err != nil {
    panic(err)
  }

  h := grpchandler.New(ctrl)
  m, err := membership.New(membership.Config{
    NodeName: cfg.NodeName,
    BindAddr: cfg.SerfAddr,
    Tags: map[string]string{
      "raft_addr": cfg.RaftAddr,
    },
    SerfJoinAddrs: cfg.SerfJoinAddrs,
  }, ctrl)
  if err != nil {
    panic(err)
  }

  srv := grpc.NewServer()
  api.RegisterMessagesServer(srv, h)

  a := &GRPCMessagesServer{
    ln: ln,
    srv: srv,
    m: m,
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