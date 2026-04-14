package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/files"
	"github.com/bd878/gallery/server/files/config"
	"github.com/bd878/gallery/server/messages/migrations"
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
		RpcAddr:               cfg.RpcAddr,
		NodeName:              cfg.NodeName,
		LogLevel:              cfg.LogLevel,
		NatsAddr:              cfg.NatsAddr,
		PGConn:                cfg.PGConn,
		GooseTableName:        cfg.GooseTableName,
	})
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		if err1 := db.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to close db", err1)
		}
	}(s.DB())

	err = s.MigrateDB(migrations.FS)
	if err != nil {
		panic(err)
	}

	err = files.Root(s.Waiter().Context(), cfg, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
		panic(err)
	}

	fmt.Println("starting files service")
	defer fmt.Println("stopped files service")

	s.Waiter().Add(
		s.WaitForPool,
		s.WaitForStream,
		s.WaitForMux,
		s.WaitForRPC,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}