package main

import (
  "flag"
  "encoding/json"
  "os"
  "log"

  "github.com/bd878/gallery/server/messages/config"
)

var (
  configPath = flag.String("config", "config/default.json", "config path")
  interactive = flag.Bool("interactive", true, "ignore logFile in config " + 
    "output log messages to stdout")
)

func main() {
  flag.Parse()

  c := loadConfig()

  if c.Debug {
    if *interactive {
      log.SetOutput(os.Stdout)
    } else {
      f := setLogOutput(c.LogFile)
      defer f.Close()
    }
  }

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

func setLogOutput(p string) *os.File {
  f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }

  log.SetOutput(f)
  return f
}
