package main

import (
  "flag"
  "fmt"
  "path/filepath"
  "encoding/json"
  "os"
  "log"

  "github.com/bd878/gallery/server/messages/config"
)

var (
  configPath = flag.String("config", "config/default.json", "config path")
)

func main() {
  flag.Parse()

  c := loadConfig()

  f := setLogOutput(c.LogPath, c.NodeName)
  defer f.Close()

  fmt.Printf("=== GRPC %s\n", c.NodeName)
  fmt.Println("Addr:", c.RpcAddr)
  fmt.Println("LogFile:", f.Name())
  fmt.Println()

  log.Printf("=== GRPC %s\n", c.NodeName)

  server := New(c)
  server.Run()
}

func loadConfig() config.Config {
  f, err := os.Open(*configPath)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  var cfg config.Config
  if err := json.NewDecoder(f).Decode(&cfg); err != nil {
    panic(err)
  }

  return cfg
}

func setLogOutput(dir, nodeName string) *os.File {
  if err := os.MkdirAll(dir, 0750); err != nil {
    panic(err)
  }

  logFile := fmt.Sprintf("%s.log", nodeName)

  f, err := os.OpenFile(
    filepath.Join(dir, logFile),
    os.O_APPEND|os.O_CREATE|os.O_WRONLY,
    0644,
  )
  if err != nil {
    panic(err)
  }

  log.SetOutput(f)
  return f
}
