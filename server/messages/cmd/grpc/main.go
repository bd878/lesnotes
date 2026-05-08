package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/messages"
	"github.com/bd878/gallery/server/messages/config"
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
		RpcAddr:        cfg.RpcAddr,
		NodeName:       cfg.NodeName,
		LogLevel:       cfg.LogLevel,
		SkipCaller:     1,
		NatsAddr:       cfg.NatsAddr,
		PGConn:         cfg.PGConn,
		GooseTableName: cfg.GooseTableName,
	})
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		if err1 := s.ResetDB(); err1 != nil {
			fmt.Fprintln(os.Stderr, "failed to reset db", err1)
		}

		if err2 := db.Close(); err2 != nil {
			fmt.Fprintln(os.Stderr, "failed to close db", err2)
		}
	}(s.DB())

	err = s.MigrateDB(migrations.FS)
	if err != nil {
		panic(err)
	}

	err = messages.Root(s.Waiter().Context(), cfg, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server exited %v\n", err)
		panic(err)
	}

	fmt.Println("starting messages service")
	defer fmt.Println("stopped messages service")

	s.Waiter().Add(
		s.WaitForPool,
		s.WaitForStream,
		s.WaitForMux,
		s.WaitForRPC,
		s.WaitForChannelz,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}
