package messages

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/discovery/serf"
	"github.com/bd878/gallery/server/internal/consensus/raft"
	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/messages/internal/repository/postgres"
	"github.com/bd878/gallery/server/messages/internal/machine"
	"github.com/bd878/gallery/server/messages/internal/handler/stream"
	"github.com/bd878/gallery/server/messages/internal/controller/distributed"
	"github.com/bd878/gallery/server/messages/internal/handler/grpc"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	messagesRepo := postgres.NewMessagesRepository(svc.Pool(), "messages.messages")
	filesRepo := postgres.NewFilesRepository(svc.Pool(), "messages.files")
	translationsRepo := postgres.NewTranslationsRepository(svc.Pool(), "messages.translations")

	consensus, err := setupRaft(svc, cfg, messagesRepo, filesRepo, translationsRepo)
	if err != nil {
		return err
	}

	if err := setupSerf(svc.Waiter().Context(), cfg, consensus, svc.Logger()); err != nil {
		return err
	}

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher, stream.NewDomainEventHandlers(nats.NewStream(svc.Nats())))

	controller := application.New(consensus, dispatcher, messagesRepo, filesRepo, translationsRepo, svc.Logger())

	messagesHandler := grpc.NewMessagesHandler(controller)
	translationsHandler := grpc.NewTranslationsHandler(controller)

	api.RegisterMessagesServer(svc.RPC(), messagesHandler)
	api.RegisterTranslationsServer(svc.RPC(), translationsHandler)

	return nil
}

func setupSerf(ctx context.Context, cfg config.Config, handler serf.Handler, logger *logger.Logger) error {
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

	go func() {
		defer func() {
			if err := membership.Leave(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()

		membership.Run(ctx)
	}()

	return nil
}

func setupRaft(svc system.Service, cfg config.Config, messagesRepo *postgres.MessagesRepository,
	filesRepo *postgres.FilesRepository, translationsRepo *postgres.TranslationsRepository) (*raft.Distributed, error) {
	raftListener := svc.Mux().Match(func(r io.Reader) bool {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(raft.RaftRPC)}) == 0
	})

	fsm := machine.New(messagesRepo, filesRepo, translationsRepo, svc.Logger())

	consensus, err := raft.New(raft.Config{
		Bootstrap:      cfg.RaftBootstrap,
		NodeName:       cfg.NodeName,
		RaftLogLevel:   cfg.RaftLogLevel,
		DataDir:        cfg.DataPath,
		Servers:        cfg.RaftServers,
	}, raft.NewStreamLayer(raftListener), fsm, svc.Logger())
	if err != nil {
		return nil, err
	}

	svc.Waiter().Add(func(ctx context.Context) error {
		<-ctx.Done()
		return consensus.Leave(consensus.NodeName())
	})

	return consensus, nil
}
