package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"github.com/bd878/gallery/server/users/internal/grpc"
	"github.com/bd878/gallery/server/users/config"
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
		LogPath:          cfg.LogPath,
		NodeName:         cfg.NodeName,
		SkipCaller:       0,
	}))

	server := grpc.New(grpc.Config{
		Addr:                 cfg.RpcAddr,
		PGConn:               cfg.PGConn,
		NodeName:             cfg.NodeName,
		DataPath:             cfg.DataPath,
		SessionsServiceAddr:  cfg.SessionsServiceAddr,
		MessagesServiceAddr:  cfg.MessagesServiceAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
