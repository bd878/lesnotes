package main

import (
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

  fcgihandler "github.com/bd878/gallery/server/internal/handler/fcgi"
  controller "github.com/bd878/gallery/server/internal/controller/messages"
  memory "github.com/bd878/gallery/server/internal/repository/memory"
)

func main() {
  serverCfg := loadConfig()

  if serverCfg.Debug {
    f := setLogOutput(serverCfg.LogFile)
    defer f.Close() 
  }

  c := make(chan os.Signal, 1)
  go trackConfig(c)
  defer close(c)

  mem := memory.New()
  ctrl := controller.New(mem)
  h := fcgihandler.New(ctrl)

  netCfg := net.ListenConfig{}
  l, err := netCfg.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", serverCfg.Port))
  if err != nil {
    panic(err)
  }
  defer l.Close()

  http.Handle("/send_message", http.HandlerFunc(h.SaveMessage))
  http.Handle("/read_messages", http.HandlerFunc(h.ReadMessages))
  http.Handle("/status", http.HandlerFunc(h.ReportStatus))

  log.Println("server is listening on =", serverCfg.Port)
  if err := fcgi.Serve(l, nil); err != nil {
    panic(err)
  }
}

func loadConfig() *config {
  f, err := os.Open("base.json")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  var cfg config
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
        f = setLogOutput(cfg.LogFile)
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
