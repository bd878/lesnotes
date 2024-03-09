package agent

import (
  "context"
  "net"
  "net/http"
  "net/http/fcgi"

  discovery "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  fcgihandler "github.com/bd878/gallery/server/messages/internal/handler/fcgi"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  messagesCtrl "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  sqlite "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
)

type Agent struct {
  Config

  ch chan struct{}
  listener net.Listener
  membership *discovery.Membership
}

type Config struct {
  UserAddr string
  BindAddr string
  DiscoveryAddr string 
  DBPath string
  DataPath string

  Bootstrap bool
  NodeName string
  StartJoinAddrs []string
}

func New(config Config) (*Agent, error) {
  a := &Agent{
    Config: config,
    ch: make(chan struct{}, 1),
  }

  if err := a.setupServer(); err != nil {
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
  if err := fcgi.Serve(a.listener, nil); err != nil {
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
    BindAddr: a.Config.BindAddr,
    Tags: map[string]string{
      "rpc_addr": a.Config.DiscoveryAddr,
    },
    StartJoinAddrs: a.Config.StartJoinAddrs,
  })
  return err
}

func (a *Agent) setupServer() error {
  mem, err := sqlite.New(a.Config.DBPath)
  if err != nil {
    return err
  }
  messagesConfig := messagesCtrl.Config{}
  msgCtrl, err := messagesCtrl.New(mem, messagesConfig)
  if err != nil {
    return err
  }
  userGateway := usergateway.New(a.Config.UserAddr)
  h := fcgihandler.New(msgCtrl, userGateway, a.Config.DataPath)

  netCfg := net.ListenConfig{}
  l, err := netCfg.Listen(context.Background(), "tcp4", a.Config.BindAddr)
  if err != nil {
    return err
  }
  a.listener = l

  http.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SaveMessage)))
  http.Handle("/messages/v1/read", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  http.Handle("/messages/v1/status", http.HandlerFunc(h.ReportStatus))

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