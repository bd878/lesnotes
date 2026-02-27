package search

import (
	"os"
	"fmt"
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/discovery/serf"
	"github.com/bd878/gallery/server/internal/consensus/raft"
	"github.com/bd878/gallery/server/search/config"
	"github.com/bd878/gallery/server/search/internal/repository/postgres"
	"github.com/bd878/gallery/server/search/internal/machine"
	"github.com/bd878/gallery/server/search/internal/handler/stream"
	"github.com/bd878/gallery/server/search/internal/controller/distributed"
	"github.com/bd878/gallery/server/search/internal/handler/grpc"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	messagesRepo := postgres.NewMessagesRepository(svc.Pool(), "search.messages")
	filesRepo := postgres.NewFilesRepository(svc.Pool(), "search.files")
	threadsRepo := postgres.NewThreadsRepository(svc.Pool(), "search.threads")
	translationsRepo := postgres.NewTranslationsRepository(svc.Pool(), "search.translations")

	consensus, err := setupRaft(svc, cfg, messagesRepo, threadsRepo, filesRepo, translationsRepo)
	if err != nil {
		return err
	}

	if err := setupSerf(svc, cfg, consensus, svc.Logger()); err != nil {
		return err
	}

	controller := application.New(consensus, messagesRepo, filesRepo, threadsRepo, translationsRepo, svc.Logger())

	stream.RegisterIntegrationEventHandlers(nats.NewStream(svc.Nats()),
		stream.NewIntegrationEventHandlers(controller, controller, controller, controller, svc.Logger()))

	handler := grpc.New(controller)

	api.RegisterSearchServer(svc.RPC(), handler)
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

func setupRaft(svc system.Service, cfg config.Config, messagesRepo *postgres.MessagesRepository, threadsRepo *postgres.ThreadsRepository,
	filesRepo *postgres.FilesRepository, translationsRepo *postgres.TranslationsRepository) (*raft.Distributed, error) {

	fsm := machine.New(messagesRepo, filesRepo, threadsRepo, translationsRepo, svc.Logger())

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
