package main

import (
  "fmt"
  "net"
  "google.golang.org/grpc"
  "github.com/hashicorp/raft"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/config"

  repository "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  controller "github.com/bd878/gallery/server/messages/internal/controller/distributed"
  grpchandler "github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

type Agent struct {
  ln   net.Listener
  srv *grpc.Server
}

func New(cfg config.Config) *Agent {
  ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", cfg.Port))
  if err != nil {
    panic(err)
  }

  streamLn, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.StreamPort))
  if err != nil {
    panic(err)
  }

  repo, err := repository.New(cfg.DBPath)
  if err != nil {
    panic(err)
  }

  ctrl, _ := controller.New(repo, controller.Config{
    Raft: raft.Config{
      LocalID: raft.ServerID(fmt.Sprintf("%d", cfg.StreamPort)),
    },
    StreamLayer: controller.NewStreamLayer(streamLn),
    Bootstrap:   cfg.Bootstrap,
    DataDir:     cfg.DataPath,
  })

  h := grpchandler.New(ctrl)

  srv := grpc.NewServer()
  api.RegisterMessagesServer(srv, h)

  a := &Agent{
    ln: ln,
    srv: srv,
  }

  return a
}

func (a *Agent) Run() {
  defer a.ln.Close()
  a.srv.Serve(a.ln)
}