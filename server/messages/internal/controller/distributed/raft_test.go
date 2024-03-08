package messages_test

import (
  "testing"
  "net"
  sqlite "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
  distributed "github.com/bd878/gallery/server/messages/internal/controller/distributed"
)

func TestDistributed(t *testing.T) {
  mem, err := sqlite.New(t.TempDir())
  if err != nil {
    panic(err)
  }

  ln, err := net.Listen("tcp", "127.0.0.1:8080")
  if err != nil {
    panic(err)
  }

  streamLayer := distributed.NewStreamLayer(ln)
  config := &distributed.Config{
    Bootstrap: true,
    DataDir: t.TempDir(),
    StreamLayer: streamLayer,
  }
  _, err = distributed.NewDistributedMessages(mem, config)
  if err != nil {
    t.Error(err)
  }
}