package main

import (
  "flag"
  "context"
  "io"
  "syscall"
  "encoding/json"
  "os"
  "net"
  "net/http"
  "net/http/fcgi"
  "log"
  "fmt"
  "os/signal"

  configs "github.com/bd878/gallery/server/messages/configs"
  fcgihandler "github.com/bd878/gallery/server/messages/internal/handler/fcgi"
  usergateway "github.com/bd878/gallery/server/messages/internal/gateway/user/grpc"
  messagesCtrl "github.com/bd878/gallery/server/messages/internal/controller/messages"
  sqlite "github.com/bd878/gallery/server/messages/internal/repository/sqlite"
)

var (
  configPath = flag.String("config", "base.json", "config path")
  interactive = flag.Bool("interactive", false, "ignore logFile in config " + 
    "output log messages to stdout")
)

func main() {
  flag.Parse()

  serverCfg := loadConfig()

  if serverCfg.Debug {
    if *interactive {
      log.SetOutput(os.Stdout)
    } else {
      f := setLogOutput(serverCfg.LogFile)
      defer f.Close()
    }
  }

  c := make(chan os.Signal, 1)
  go trackConfig(c)
  defer close(c)

  mem, err := sqlite.New(serverCfg.DBPath)
  if err != nil {
    panic(err)
  }
  msgCtrl := messagesCtrl.New(mem)
  userGateway := usergateway.New(serverCfg.UserAddr)
  h := fcgihandler.New(msgCtrl, userGateway, serverCfg.DataPath)

  netCfg := net.ListenConfig{}
  l, err := netCfg.Listen(context.Background(), "tcp4", fmt.Sprintf(":%d", serverCfg.Port))
  if err != nil {
    panic(err)
  }
  defer l.Close()

  http.Handle("/messages/v1/send", http.HandlerFunc(h.CheckAuth(h.SaveMessage)))
  http.Handle("/messages/v1/read_all", http.HandlerFunc(h.CheckAuth(h.ReadMessages)))
  http.Handle("/messages/v1/status", http.HandlerFunc(h.ReportStatus))

  log.Println("server is listening on =", l.Addr())
  if err := fcgi.Serve(l, nil); err != nil {
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

func trackConfig(c chan os.Signal) {
  signal.Notify(c, syscall.SIGHUP)

  var f *os.File
  defer f.Close()

  for {
    switch <-c {
    case syscall.SIGHUP:
      log.Println("recieve sighup")

      cfg := loadConfig()
      if cfg.Debug {
        if *interactive {
          log.SetOutput(os.Stdout)
        } else {
          f = setLogOutput(cfg.LogFile)
        }
      } else {
        f.Close()
        log.SetOutput(io.Discard)
      }
    }
  }
}

func setLogOutput(p string) *os.File {
  f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }

  log.SetOutput(f)
  return f
}
