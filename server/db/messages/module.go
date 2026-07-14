package messages

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/internal/consensus/raft"
	"github.com/bd878/gallery/server/internal/discovery/serf"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/db/messages/config"
	"github.com/bd878/gallery/server/db/messages/internal/controller/distributed"
	"github.com/bd878/gallery/server/db/messages/internal/handler/grpc"
	"github.com/bd878/gallery/server/db/messages/internal/machine"
	"github.com/bd878/gallery/server/db/messages/internal/repository/postgres"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	messagesRepo := postgres.NewMessagesRepository(svc.Pool(), "messages.messages")
	translationsRepo := postgres.NewTranslationsRepository(svc.Pool(), "messages.translations")
	commentsRepo := postgres.NewCommentsRepository(svc.Pool(), "messages.comments")
	dumper := postgres.NewDumper(svc.Pool(), "messages.messages", "messages.translations", "messages.comments")

	consensus, err := setupRaft(svc, cfg, messagesRepo, translationsRepo, commentsRepo, dumper)
	if err != nil {
		return err
	}

	if err := setupSerf(svc, cfg, consensus, svc.Logger()); err != nil {
		return err
	}

	controller := application.New(consensus, messagesRepo,
			translationsRepo, commentsRepo, svc.Logger())

	messagesHandler := grpc.NewMessagesHandler(controller)
	translationsHandler := grpc.NewTranslationsHandler(controller)
	commentsHandler := grpc.NewCommentsHandler(controller)

	messages.RegisterMessagesServer(svc.RPC(), messagesHandler)
	api.RegisterDistributedServer(svc.RPC(), messagesHandler)
	translations.RegisterTranslationsServer(svc.RPC(), translationsHandler)
	comments.RegisterCommentsServer(svc.RPC(), commentsHandler)

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

func setupRaft(svc system.Service, cfg config.Config, messagesRepo *postgres.MessagesRepository,
	translationsRepo *postgres.TranslationsRepository, commentsRepo *postgres.CommentsRepository, dumper *postgres.Dumper) (*raft.Distributed, error) {
	fsm := machine.New(messagesRepo, translationsRepo, commentsRepo, dumper, svc.Logger())

	consensus, err := raft.New(raft.Config{
		Bootstrap:    cfg.RaftBootstrap,
		NodeName:     cfg.NodeName,
		RaftLogLevel: cfg.RaftLogLevel,
		DataDir:      cfg.DataPath,
		Servers:      cfg.RaftServers,
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
