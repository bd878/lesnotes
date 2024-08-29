package main

import (
  "net"
  "io"
  "bytes"
  "google.golang.org/grpc"
  "github.com/hashicorp/raft"
  "github.com/soheilhy/cmux"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/config"
  hclog "github.com/hashicorp/go-hclog"

  membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type GRPCMessagesServer struct {
  cfg     config.Config
  mux     cmux.CMux
  ln      net.Listener
  server *grpc.Server
  ctrl   *controller.DistributedMessages
  m      *membership.Membership
}

func New(cfg config.Config) *GRPCMessagesServer {
  ln, err := net.Listen("tcp4", cfg.RpcAddr)
  if err != nil {
    panic(err)
  }

  mux := cmux.New(ln)

  s := &GRPCMessagesServer{
    cfg:  cfg,
    mux:  mux,
    ln:   ln,
  }

  s.setupRaft()
  s.setupGRPC()

  return s
}

func (s *GRPCMessagesServer) setupRaft() {
  repo, err := repository.New(s.cfg.DBPath)
  if err != nil {
    panic(err)
  }

  raftLogLevel := hclog.Error.String()
  switch s.cfg.RaftLogLevel {
  case "debug":
    raftLogLevel = hclog.Debug.String()
  case "error":
    raftLogLevel = hclog.Error.String()
  case "info":
    raftLogLevel = hclog.Info.String()
  }

  raftLn := s.mux.Match(func(r io.Reader) bool {
    b := make([]byte, 1)
    if _, err := r.Read(b); err != nil {
      return false
    }
    return bytes.Compare(b, []byte{byte(controller.RaftRPC)}) == 0
  })

  s.ctrl, err = controller.New(repo, controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(s.cfg.NodeName),
      LogLevel: raftLogLevel,
    },
    StreamLayer: controller.NewStreamLayer(raftLn),
    Bootstrap:   s.cfg.RaftBootstrap,
    DataDir:     s.cfg.DataPath,
    Servers:     s.cfg.RaftServers,
  })
  if err != nil {
    panic(err)
  }
}

func (s *GRPCMessagesServer) setupGRPC() {
  var err error
  h := grpchandler.New(s.ctrl)
  s.m, err = membership.New(
    membership.Config{
      NodeName: s.cfg.NodeName,
      BindAddr: s.cfg.SerfAddr,
      Tags: map[string]string{
        "raft_addr": s.cfg.RpcAddr,
      },
      SerfJoinAddrs: s.cfg.SerfJoinAddrs,
    },
    s.ctrl,
  )
  if err != nil {
    panic(err)
  }

  s.server = grpc.NewServer()
  // TODO: MessagesSerivce &api.MessagesService{
  //    Produce: s.server.Produce,
  //    Consume: s.server.Consume,
  //    ...
  // }
  api.RegisterMessagesServer(s.server, h)

  grpcLn := s.mux.Match(cmux.Any())
  go func() {
    if err := s.server.Serve(grpcLn); err != nil {
      s.Shutdown()
    }
  }()
}

func (s *GRPCMessagesServer) Run() {
  defer s.mux.Close()
  s.mux.Serve()
}

func (s *GRPCMessagesServer) Addr() string {
  return s.ln.Addr().String()
}

func (s *GRPCMessagesServer) Shutdown() {
/* TODO: implement, Leave server from cluster */
}