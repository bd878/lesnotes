package main

import (
	"flag"
	"fmt"
	"context"
	"os"

	"github.com/bd878/gallery/server/users/config"
	"github.com/bd878/gallery/server/users/internal/http"
	"github.com/bd878/gallery/server/logger"
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
		NodeName:         cfg.NodeName,
		LogLevel:         cfg.LogLevel,
		SkipCaller:       0,
	}))

	server := http.New(http.Config{
		Addr:                cfg.HttpAddr,
		RpcAddr:             cfg.RpcAddr,
		CookieDomain:        cfg.CookieDomain,
		MessagesServiceAddr: cfg.MessagesServiceAddr,
		SessionsServiceAddr: cfg.SessionsServiceAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
