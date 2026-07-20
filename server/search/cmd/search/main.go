package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/bd878/gallery/server/db/search/pkg/loadbalance"
	_ "github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"

	"github.com/bd878/gallery/server/search/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/search"
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
	})
	if err != nil {
		panic(err)
	}

	if err := search.Root(s.Waiter().Context(), cfg, s); err != nil {
		panic(err)
	}

	fmt.Println("starting search service")
	defer fmt.Println("stopped search service")

	s.Waiter().Add(
		s.WaitForHTTP,
		s.WaitForStream,
	)

	if err = s.Waiter().Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "waiter exited with error", err)
	}
}
