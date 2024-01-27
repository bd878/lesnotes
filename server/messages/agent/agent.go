package agent

import (
  "context"
  "net"
  "net/http"
  "net/http/fcgi"

  // membership "github.com/bd878/gallery/server/messages/internal/discovery/serf"
  fcgihandler "github.com/bd878/gallery/server/messages/internal/handler/fcgi"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  messagesCtrl "github.com/bd878/gallery/server/messages/internal/controller/messages"
  sqlite "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
)

type Agent struct {
  Config

  listener net.Listener
}

type Config struct {
  UserAddr string
  BindAddr string
  DBPath string
  DataPath string

  Bootstrap bool
  NodeName string
  StartJoinAddrs []string
}

func New(config Config) (*Agent, error) {
  a := &Agent{Config: config}

  if err := a.setupServer(); err != nil {
    return nil, err
  }

  return a, nil
}

func (a *Agent) Serve() error {
  return a.serve()
}

func (a *Agent) setupServer() error {
  mem, err := sqlite.New(a.Config.DBPath)
  if err != nil {
    panic(err)
  }
  msgCtrl := messagesCtrl.New(mem)
  userGateway := usergateway.New(a.Config.UserAddr)
  h := fcgihandler.New(msgCtrl, userGateway, a.Config.DataPath)

  netCfg := net.ListenConfig{}
  l, err := netCfg.Listen(context.Background(), "tcp4", a.Config.BindAddr)
  if err != nil {
    panic(err)
  }
  a.listener = l

  http.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SaveMessage)))
  http.Handle("/messages/v1/read", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  http.Handle("/messages/v1/status", http.HandlerFunc(h.ReportStatus))

  return nil
}

func (a *Agent) serve() error {
  if err := fcgi.Serve(a.listener, nil); err != nil {
    return err
  }
  return nil
}

func (a *Agent) Shutdown() error {
  if err := a.listener.Close(); err != nil {
    return err
  }
  return nil
}
