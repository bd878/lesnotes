package main

import (
  "net"
  "io"
  "sync"
  "context"
  "bytes"
  "google.golang.org/grpc"
  "github.com/hashicorp/raft"
  "github.com/soheilhy/cmux"

  "github.com/bd878/gallery/server/api"
  hclog "github.com/hashicorp/go-hclog"

  membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type GRPCServerConfig struct {
  Addr            string
  DBPath          string
  NodeName        string
  RaftLogLevel    string
  RaftBootstrap   bool
  DataPath        string
  RaftServers     []string
  SerfAddr        string
  SerfJoinAddrs   []string
}

type GRPCServer struct {
  config  GRPCServerConfig
  mux     cmux.CMux
  ln      net.Listener
  server *grpc.Server
  ctrl   *controller.DistributedMessages
  m      *membership.Membership
}

func NewGRPCServer(cfg GRPCServerConfig) *GRPCServer {
  ln, err := net.Listen("tcp4", cfg.Addr)
  if err != nil {
    panic(err)
  }

  mux := cmux.New(ln)

  s := &GRPCServer{
    config:  cfg,
    mux:     mux,
    ln:      ln,
  }

  s.setupRaft()
  s.setupGRPC()

  return s
}

func (s *GRPCServer) setupRaft() {
  repo, err := repository.New(s.config.DBPath)
  if err != nil {
    panic(err)
  }

  raftLogLevel := hclog.Error.String()
  switch s.config.RaftLogLevel {
  case "debug":
    raftLogLevel = hclog.Debug.String()
  case "error":
    raftLogLevel = hclog.Error.String()
  case "info":
    raftLogLevel = hclog.Info.String()
  default:
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
      LocalID: raft.ServerID(s.config.NodeName),
      LogLevel: raftLogLevel,
    },
    StreamLayer: controller.NewStreamLayer(raftLn),
    Bootstrap:   s.config.RaftBootstrap,
    DataDir:     s.config.DataPath,
    Servers:     s.config.RaftServers,
  })

  if err != nil {
    panic(err)
  }
}

func (s *GRPCServer) setupGRPC() {
  var err error
  h := grpchandler.New(s.ctrl)
  s.m, err = membership.New(
    membership.Config{
      NodeName: s.config.NodeName,
      BindAddr: s.config.SerfAddr,
      Tags: map[string]string{
        "raft_addr": s.config.Addr,
      },
      SerfJoinAddrs: s.config.SerfJoinAddrs,
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

func (s *GRPCServer) Run(ctx context.Context, wg *sync.WaitGroup) {
  defer s.mux.Close()
  s.mux.Serve()
  wg.Done()
}

func (s *GRPCServer) Addr() string {
  return s.ln.Addr().String()
}

func (s *GRPCServer) Shutdown() {
/* TODO: implement, Leave server from cluster */
}