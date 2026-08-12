package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bd878/gallery/server/threads/migrations"
	_ "github.com/bd878/gallery/server/db/threads/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"

	"github.com/bd878/gallery/server/threads/config"
	"github.com/bd878/gallery/server/threads"
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
		NodeName: cfg.NodeName,
		LogLevel: cfg.LogLevel,
		SkipCaller: 1,
		NatsStream: cfg.NatsStream,
		NatsAddr: cfg.NatsAddr,
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

	if err := threads.Root(s.Waiter().Context(), cfg, s); err != nil {
		panic(err)
	}

	fmt.Println("starting threads service")
	defer fmt.Println("stopped threads service")

	s.Waiter().Add(
		s.WaitForHTTP,
		s.WaitForStream,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}
