package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"github.com/bd878/gallery/server/messages/internal/grpc"
	"github.com/bd878/gallery/server/messages/config"
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
		LogPath:    cfg.LogPath,
		NodeName:   cfg.NodeName,
		LogLevel:   cfg.LogLevel,
		SkipCaller: 0,
	}))

	server := grpc.New(grpc.Config{
		Addr:                  cfg.RpcAddr,
		PGConn:                cfg.PGConn,
		NodeName:              cfg.NodeName,
		RaftLogLevel:          cfg.RaftLogLevel,
		RaftBootstrap:         cfg.RaftBootstrap,
		DataPath:              cfg.DataPath,
		RaftServers:           cfg.RaftServers,
		SerfAddr:              cfg.SerfAddr,
		SerfJoinAddrs:         cfg.SerfJoinAddrs,
		SessionsServiceAddr:   cfg.SessionsServiceAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
