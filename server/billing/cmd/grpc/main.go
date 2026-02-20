package main

import (
	"flag"
	"fmt"
	"os"
	"database/sql"

	"github.com/bd878/gallery/server/billing"
	"github.com/bd878/gallery/server/billing/migrations"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/internal/system"
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
		RpcAddr:            cfg.RpcAddr,
		NodeName:           cfg.NodeName,
		LogLevel:           cfg.LogLevel,
		SkipCaller:         1,
		NatsAddr:           cfg.NatsAddr,
		PGConn:             cfg.PGConn,
		GooseTableName:     cfg.GooseTableName,
	})
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := s.ResetDB()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to reset db", err)
		}

		err = db.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to close db", err)
		}
	}(s.DB())

	if err := s.MigrateDB(migrations.FS); err != nil {
		panic(err)
	}

	if err := billing.Root(s.Waiter().Context(), cfg, s); err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
		panic(err)
	}

	fmt.Println("started billing service")
	defer fmt.Println("stopped billing service")

	s.Waiter().Add(
		s.WaitForPool,
		s.WaitForStream,
		s.WaitForRPC,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}
