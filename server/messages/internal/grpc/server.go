package grpc

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
  "github.com/bd878/gallery/server/logger"
  hclog "github.com/hashicorp/go-hclog"

  grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
  membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type Config struct {
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

type Server struct {
  *grpc.Server
  conf             Config
  mux              cmux.CMux
  listener         net.Listener
  grpcListener     net.Listener
  controller      *controller.DistributedMessages
  membership      *membership.Membership
}

func New(cfg Config) *Server {
  listener, err := net.Listen("tcp4", cfg.Addr)
  if err != nil {
    panic(err)
  }

  mux := cmux.New(listener)

  server := &Server{
    conf:          cfg,
    mux:           mux,
    listener:      listener,
  }

  server.setupRaft(logger.Default())
  server.setupGRPC(logger.Default())

  return server
}

func (s *Server) setupRaft(log *logger.Logger) {
  repo, err := repository.New(s.conf.DBPath)
  if err != nil {
    panic(err)
  }

  raftLogLevel := hclog.Error.String()
  switch s.conf.RaftLogLevel {
  case "debug":
    raftLogLevel = hclog.Debug.String()
  case "error":
    raftLogLevel = hclog.Error.String()
  case "info":
    raftLogLevel = hclog.Info.String()
  default:
    raftLogLevel = hclog.Info.String()
  }

  raftListener := s.mux.Match(func(r io.Reader) bool {
    b := make([]byte, 1)
    if _, err := r.Read(b); err != nil {
      return false
    }
    return bytes.Compare(b, []byte{byte(controller.RaftRPC)}) == 0
  })

  control, err := controller.New(repo, controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(s.conf.NodeName),
      LogLevel: raftLogLevel,
    },
    StreamLayer: controller.NewStreamLayer(raftListener),
    Bootstrap:   s.conf.RaftBootstrap,
    DataDir:     s.conf.DataPath,
    Servers:     s.conf.RaftServers,
  })
  if err != nil {
    panic(err)
  }

  s.controller = control
}

func (s *Server) setupGRPC(log *logger.Logger) {
  handler := grpchandler.New(s.controller)
  member, err := membership.New(
    membership.Config{
      NodeName: s.conf.NodeName,
      BindAddr: s.conf.SerfAddr,
      Tags: map[string]string{
        "raft_addr": s.conf.Addr,
      },
      SerfJoinAddrs: s.conf.SerfJoinAddrs,
    },
    s.controller,
  )
  if err != nil {
    log.Error("message", "failed to establish membership connection")
    panic(err)
  }

  s.membership = member

  s.Server = grpc.NewServer(
    grpc.ChainUnaryInterceptor(
      grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
    ),
  )

  api.RegisterMessagesServer(s.Server, handler)

  grpcListener := s.mux.Match(cmux.Any())
  s.grpcListener = grpcListener
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
  go s.Serve(s.grpcListener)
  defer s.mux.Close()
  s.mux.Serve()
  wg.Done()
}

func (s *Server) Addr() string {
  return s.listener.Addr().String()
}

func (s *Server) Shutdown() {
/* TODO: implement, Leave server from cluster */
}