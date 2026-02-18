package main

import (
	"flag"
	"fmt"
	"os"
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/bd878/gallery/server/billing/migrations"
	"github.com/bd878/gallery/server/billing/internal/grpc"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/internal/system"
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

	s, err := system.NewSystem(system.Config{
		NodeName:        cfg.NodeName,
		LogLevel:        cfg.LogLevel,
		SkipCaller:      cfg.SkipCaller,
		NatsAddr:        cfg.NatsAddr,
		PGConn:          cfg.PGConn,
	})
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", cfg.PGConn)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := s.ResetDB()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = db.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}(s.DB())

	if err := s.MigrateDB(migrations.FS); err != nil {
		panic(err)
	}

	server := grpc.New(grpc.Config{
		Addr:                  cfg.RpcAddr,
		PGConn:                cfg.PGConn,
		NodeName:              cfg.NodeName,
		RaftLogLevel:          cfg.RaftLogLevel,
		RaftBootstrap:         cfg.RaftBootstrap,
		DataPath:              cfg.DataPath,
		PaymentsTableName:     cfg.PaymentsTableName,
		InvoicesTableName:     cfg.InvoicesTableName,
		RaftServers:           cfg.RaftServers,
		SerfAddr:              cfg.SerfAddr,
		SerfJoinAddrs:         cfg.SerfJoinAddrs,
		NatsAddr:              cfg.NatsAddr,
	})

	

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}

	s.Waiter().Add(
		s.WaitForPool,
		s.WaitForStream,
		s.WaitForRPC,
	)

	return s.Waiter().Wait()
}
