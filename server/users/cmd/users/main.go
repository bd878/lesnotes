package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bd878/gallery/server/users/migrations"
	_ "github.com/bd878/gallery/server/db/users/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/db/messages/pkg/loadbalance"

	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/users/config"
	"github.com/bd878/gallery/server/users"
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
		NodeName: cfg.NodeName,
		LogLevel: cfg.LogLevel,
		SkipCaller: 1,
		NatsAddr: cfg.NatsAddr,
		NatsStream: cfg.NatsStream,
		HttpAddr: cfg.HttpAddr,
		PGConn: cfg.PGConn,
		GooseTableName: cfg.GooseTableName,
	})
	if err != nil {
		panic(err)
	}

	if err := s.MigrateDB(migrations.FS); err != nil {
		panic(err)
	}

	if err := users.Root(s.Waiter().Context(), cfg, s); err != nil {
		panic(err)
	}

	fmt.Println("starting users http service")
	defer fmt.Println("stopped users http service")

	s.Waiter().Add(
		s.WaitForHTTP,
		s.WaitForStream,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}
