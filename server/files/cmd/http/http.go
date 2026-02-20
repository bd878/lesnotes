package main

import (
	"flag"
	"fmt"
	"context"
	"os"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/files/config"
	"github.com/bd878/gallery/server/files/internal/http"
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
		SkipCaller: 0,
	}))

	server := http.New(http.Config{
		Addr:                 cfg.HttpAddr,
		RpcAddr:              cfg.RpcAddr,
		UsersServiceAddr:     cfg.UsersServiceAddr,
		SessionsServiceAddr:  cfg.SessionsServiceAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
