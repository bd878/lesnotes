package main

import (
  "flag"
  "fmt"
  "sync"
  "context"
  "os"

  "github.com/bd878/gallery/server/messages/config"
  "github.com/bd878/gallery/server/log"
)

func init() {
  flag.Usage = func() {
    fmt.Printf("Usage: %d config\n", os.Args[0])
  }
}

func main() {
  flag.Parse()

  if flag.NArg() != 1 {
    flag.Usage()
    os.Exit(1)
  }

  cfg := config.Load(flag.Arg(0))
  log.SetDefault(log.New(log.Config{
    LogPath:   cfg.LogPath,
    NodeName:  cfg.NodeName,
  }))

  server := NewHTTPServer(HTTPServerConfig{
    Addr:              cfg.HttpAddr,
    RpcAddr:           cfg.RpcAddr,
    DataPath:          cfg.DataPath,
    UserServiceAddr:   cfg.UserServiceAddr,
  })

  var wg *sync.WaitGroup
  wg.Add(1)
  go server.ListenAndServe(context.Background(), wg)
  wg.Wait()
}
