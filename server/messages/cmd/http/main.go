package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/bd878/gallery/server/messages/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/sessions/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/threads/pkg/loadbalance"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/messages/internal/http"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s config\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.Load(flag.Arg(0))
	logger.SetDefault(logger.New(logger.Config{
		NodeName:   cfg.NodeName,
		LogLevel:   cfg.LogLevel,
		SkipCaller: 0,
	}))

	server := http.New(http.Config{
		Addr:                cfg.HttpAddr,
		RpcAddr:             cfg.RpcAddr,
		UsersServiceAddr:    cfg.UsersServiceAddr,
		FilesServiceAddr:    cfg.FilesServiceAddr,
		SessionsServiceAddr: cfg.SessionsServiceAddr,
		ThreadsServiceAddr:  cfg.ThreadsServiceAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
