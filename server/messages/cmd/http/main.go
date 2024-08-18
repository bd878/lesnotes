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

  server := New(c)

  fmt.Printf("=== HTTP %s\n", c.NodeName)
  fmt.Println("Addr:", server.Addr)
  fmt.Println("LogFile:", f.Name())
  fmt.Println()

  log.Printf("=== HTTP %s\n", c.NodeName)
  server.ListenAndServe()
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
