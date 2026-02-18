package main

import (
	"flag"
	"fmt"
	"os"
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/bd878/gallery/server/messages/migrations"
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
		NodeName:   cfg.NodeName,
		LogLevel:   cfg.LogLevel,
		SkipCaller: 0,
	}))

	db, err := sql.Open("pgx", cfg.PGConn)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		if err = goose.Reset(db, "."); err != nil {
			return
		}

		if err = db.Close(); err != nil {
			return
		}
	}(db)

	goose.SetVerbose(true)
	goose.SetTableName(cfg.GooseTableName)

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}

	server := grpc.New(grpc.Config{
		Addr:                  cfg.RpcAddr,
		PGConn:                cfg.PGConn,
		NodeName:              cfg.NodeName,
		RaftLogLevel:          cfg.RaftLogLevel,
		RaftBootstrap:         cfg.RaftBootstrap,
		DataPath:              cfg.DataPath,
		TableName:             cfg.TableName,
		FilesTableName:        cfg.FilesTableName,
		TranslationsTableName: cfg.TranslationsTableName,
		RaftServers:           cfg.RaftServers,
		SerfAddr:              cfg.SerfAddr,
		SerfJoinAddrs:         cfg.SerfJoinAddrs,
		NatsAddr:              cfg.NatsAddr,
	})

	if err := server.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
	}
}
