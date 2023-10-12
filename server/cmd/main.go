package main

import (
  "flag"
  "context"
  "net"
  "net/http"
  "net/http/fcgi"
  "log"

  fcgihandler "github.com/bd878/gallery/server/internal/handler/fcgi"
  controller "github.com/bd878/gallery/server/internal/controller/messages"
  memory "github.com/bd878/gallery/server/internal/repository/memory"
)

var (
  addr = flag.String("addr", ":8083", "fcgi addr to listen")
)

func main() {
  flag.Parse()

  mem := memory.New()
  ctrl := controller.New(mem)
  h := fcgihandler.New(ctrl)

  cfg := net.ListenConfig{}
  l, err := cfg.Listen(context.Background(), "tcp", *addr)
  if err != nil {
    panic(err)
  }
  defer l.Close()

  http.Handle("/send_message", http.HandlerFunc(h.SaveMessage))
  http.Handle("/read_messages", http.HandlerFunc(h.ReadMessages))
  log.Println("server is listening on =", *addr)
  if err := fcgi.Serve(l, nil); err != nil {
    panic(err)
  }
}
