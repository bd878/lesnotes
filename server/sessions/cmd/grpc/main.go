package main

import (
	"flag"
	"fmt"
	"os"
	"context"
	"sync"

	"github.com/bd878/gallery/server/sessions/internal/grpc"
	"github.com/bd878/gallery/server/sessions/config"
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
		Addr:             cfg.RpcAddr,
		DBPath:           cfg.DBPath,
		NodeName:         cfg.NodeName,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	logger.Infoln("server is listening on:", server.Addr())
	go server.Run(context.Background(), &wg)
	wg.Wait()
}
