package threads

import (
	"os"
	"fmt"
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/discovery/serf"
	"github.com/bd878/gallery/server/internal/consensus/raft"
	"github.com/bd878/gallery/server/threads/config"
	"github.com/bd878/gallery/server/threads/internal/repository/postgres"
	"github.com/bd878/gallery/server/threads/internal/machine"
	"github.com/bd878/gallery/server/threads/internal/handler/stream"
	"github.com/bd878/gallery/server/threads/internal/controller/distributed"
	"github.com/bd878/gallery/server/threads/internal/handler/grpc"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	threadsRepo := postgres.NewThreadsRepository(svc.Pool(), "threads.threads")

	consensus, err := setupRaft(svc, cfg, threadsRepo)
	if err != nil {
		return err
	}

	if err := setupSerf(svc, cfg, consensus, svc.Logger()); err != nil {
		return err
	}

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher, stream.NewDomainEventHandlers(nats.NewStream(svc.Nats())))

	controller := application.New(consensus, dispatcher, threadsRepo, svc.Logger())

	handler := grpc.New(controller)

	api.RegisterThreadsServer(svc.RPC(), handler)
	api.RegisterDistributedServer(svc.RPC(), handler)

	return nil
}

func setupSerf(svc system.Service, cfg config.Config, handler serf.Handler, logger *logger.Logger) error {
	membership, err := serf.New(
		serf.Config{
			NodeName: cfg.NodeName,
			BindAddr: cfg.SerfAddr,
			Tags: map[string]string{
				"raft_addr": cfg.RpcAddr,
			},
			SerfJoinAddrs: cfg.SerfJoinAddrs,
		},
		handler,
	)
	if err != nil {
		return err
	}

	svc.Waiter().Add(func(ctx context.Context) (err error) {
		group, gCtx := errgroup.WithContext(ctx)

		group.Go(func() error {
			fmt.Fprintln(os.Stdout, "membership run")
			membership.Run(gCtx)
			return nil
		})

		group.Go(func() error {
			<-gCtx.Done()
			fmt.Fprintln(os.Stdout, "membership is about to leave")
			return membership.Leave()
		})

		return group.Wait()
	})

	return nil
}

func setupRaft(svc system.Service, cfg config.Config, threadsRepo *postgres.ThreadsRepository) (*raft.Distributed, error) {
	fsm := machine.New(threadsRepo, svc.Logger())

	consensus, err := raft.New(raft.Config{
		Bootstrap:      cfg.RaftBootstrap,
		NodeName:       cfg.NodeName,
		RaftLogLevel:   cfg.RaftLogLevel,
		DataDir:        cfg.DataPath,
		Servers:        cfg.RaftServers,
	}, raft.NewStreamLayer(svc.RaftListener()), fsm, svc.Logger())
	if err != nil {
		return nil, err
	}

	svc.Waiter().Add(func(ctx context.Context) error {
		<-ctx.Done()
		fmt.Fprintln(os.Stdout, "raft is about to leave")
		return consensus.Leave(consensus.NodeName())
	})

	return consensus, nil
}
