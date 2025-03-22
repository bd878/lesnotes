package main

import (
	"flag"
	"fmt"
	"sync"
	"context"
	"os"

	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/messages/internal/http"
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
		SkipCaller: 1,
	}))

	server := http.New(http.Config{
		Addr:              cfg.HttpAddr,
		RpcAddr:           cfg.RpcAddr,
		DataPath:          cfg.DataPath,
		UsersServiceAddr:  cfg.UsersServiceAddr,
		FilesServiceAddr:  cfg.FilesServiceAddr,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	logger.Infoln("server is listening on:", server.Addr)
	go server.ListenAndServe(context.Background(), &wg)
	wg.Wait()
}
