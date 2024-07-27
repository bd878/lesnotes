package agent

import (
  "context"
  "net"
  "log"
  "net/http"
  "github.com/hashicorp/raft"

  discovery "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  sqlite "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
)

type Agent struct {
  Config

  ch chan struct{}
  listener net.Listener
  membership *discovery.Membership
}

type Config struct {
  UserAddr       string
  BindAddr       string
  StreamAddr     string
  DiscoveryAddr  string 
  DBPath         string
  DataPath       string

  Bootstrap      bool
  NodeName       string
  StartJoinAddrs []string
}

func New(config Config) (*Agent, error) {
  a := &Agent{
    Config: config,
    ch: make(chan struct{}, 1),
  }

  if err := a.setupHTTPServer(); err != nil {
    return nil, err
  }

  if err := a.setupGRPCServer(); err != nil {
    return nil, err
  }

  if err := a.setupMembership(); err != nil {
    return nil, err
  }

  return a, nil
}

type handler struct {
}

func (a *Agent) Run() error {
  log.Println("http server is listening on =", a.listener.Addr().String())
  if err := http.Serve(a.listener, nil); err != nil {
    return err
  }
  return nil
}

func (h *handler) Join(_, _ string) error {
  return nil
}

func (h *handler) Leave(_ string) error {
  return nil
}

func (a *Agent) setupMembership() error {
  var err error
  a.membership, err = discovery.New(&handler{}, discovery.Config{
    NodeName: a.Config.NodeName,
    BindAddr: a.Config.DiscoveryAddr,
    Tags: map[string]string{
      "rpc_addr": a.Config.DiscoveryAddr,
    },
    StartJoinAddrs: a.Config.StartJoinAddrs,
  })
  return err
}

func (a *Agent) setupHTTPServer() error {
  mem, err := sqlite.New(a.Config.DBPath)
  if err != nil {
    return err
  }

  streamLn, err := net.Listen("tcp", a.Config.StreamAddr)
  if err != nil {
    return err
  }

  distributedCfg := controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(a.Config.BindAddr),
    },
    StreamLayer: controller.NewStreamLayer(streamLn),
    Bootstrap: a.Config.Bootstrap,
    DataDir: a.Config.DataPath,
  }
  distributedCtrl, err := controller.New(mem, distributedCfg)
  if err != nil {
    return err
  }
  userGateway := usergateway.New(a.Config.UserAddr)
  h := httphandler.New(distributedCtrl, userGateway, a.Config.DataPath)

  netCfg := net.ListenConfig{}
  l, err := netCfg.Listen(context.Background(), "tcp4", a.Config.BindAddr)
  if err != nil {
    return err
  }
  a.listener = l

  http.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SendMessage)))
  http.Handle("/messages/v1/read", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  http.Handle("/messages/v1/status", http.HandlerFunc(h.GetStatus))
  http.Handle("/messages/v1/read_file", http.HandlerFunc(h.CheckAuth(h.ReadFile)))

  return nil
}

func (a *Agent) setupGRPCServer() error {
  /* TODO: implement;
     GRPC server is necessary for replication and distribution
  */
  return nil
}

func (a *Agent) Shutdown() error {
  defer func() {
    a.ch <- struct{}{}
    close(a.ch)
  }()

  if err := a.listener.Close(); err != nil {
    return err
  }
  return nil
}

func (a *Agent) Done() chan struct{} {
  return a.ch
}