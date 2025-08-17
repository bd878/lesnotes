package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"golang.org/x/sync/errgroup"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/messages/internal/grpc"
	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/waiter"
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

	pool, err := pgxpool.New(context.Background(), cfg.PGConn)
	if err != nil {
		panic(err)
	}

	server := grpc.New(grpc.Config{
		Addr:                  cfg.RpcAddr,
		DBPath:                cfg.DBPath,
		NodeName:              cfg.NodeName,
		RaftLogLevel:          cfg.RaftLogLevel,
		RaftBootstrap:         cfg.RaftBootstrap,
		DataPath:              cfg.DataPath,
		RaftServers:           cfg.RaftServers,
		SerfAddr:              cfg.SerfAddr,
		SerfJoinAddrs:         cfg.SerfJoinAddrs,
		SessionsServiceAddr:   cfg.SessionsServiceAddr,
	})

	waiter := waiter.New(waiter.CatchSignals())

	waitForPool := func(ctx context.Context) error {
		group, gCtx := errgroup.WithContext(ctx)

		group.Go(func() error {
			<-gCtx.Done()
			fmt.Fprintln(os.Stdout, "closing pgpool connections")
			pool.Close()
			return nil
		})

		return group.Wait()
	}

	waiter.Add(
		waitForPool,
		server.Run,
	)

	waiter.Wait()
}
