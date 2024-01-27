package main

import (
  "flag"
  "encoding/json"
  "os"
  "log"
  "fmt"

  configs "github.com/bd878/gallery/server/messages/configs"
  agent "github.com/bd878/gallery/server/messages/agent"
)

var (
  configPath = flag.String("config", "base.json", "config path")
  interactive = flag.Bool("interactive", false, "ignore logFile in config " + 
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

  a, err := agent.New(agent.Config{
    UserAddr: c.UserAddr,
    BindAddr: fmt.Sprintf(":%d", c.Port),
    DBPath: c.DBPath,
    DataPath: c.DataPath,

    Bootstrap: false,
    NodeName: "messages",
    StartJoinAddrs: []string{},
  })
  if err != nil {
    panic(err)
  }

  if err := a.Serve(); err != nil {
    panic(err)
  }
}

func loadConfig() *configs.Config {
  f, err := os.Open(*configPath)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  var cfg configs.Config
  if err := json.NewDecoder(f).Decode(&cfg); err != nil {
    panic(err)
  }

  return &cfg
}

func setLogOutput(p string) *os.File {
  f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }

  log.SetOutput(f)
  return f
}
