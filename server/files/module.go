package files

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api/files"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/system"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/files/config"
	"github.com/bd878/gallery/server/files/internal/repository/postgres"
	"github.com/bd878/gallery/server/files/internal/controller/application"
	"github.com/bd878/gallery/server/files/internal/handler/grpc"
	"github.com/bd878/gallery/server/files/internal/handler/stream"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	container := di.New()

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return svc.Pool(), nil
	})
	container.AddSingleton("js", func(c di.Container) (any, error) {
		js := jetstream.NewStream(svc.Config().NatsStream, svc.JS())
		return js, nil
	})
	container.AddScoped("tx", func(c di.Container) (any, error) {
		pool := c.Get("db").(*pgxpool.Pool)
		return pool.BeginTx(ctx, pgx.TxOptions{})
	})
	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		inboxStore := pg.NewInboxStore(tx, "files_stream.inbox")
		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	filesRepo := postgres.NewFilesRepository(svc.Pool(), "files.files")
	messagesRepo := postgres.NewMessagesRepository(svc.Pool(), "files.messages")

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()

	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(am.NewEventStream(
			container.Get("js").(am.RawMessageStream),
		)),
	)

	controller := application.New(dispatcher, filesRepo, messagesRepo)

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return stream.NewIntegrationEventHandlers(controller), nil
	})

	stream.RegisterIntegrationEventHandlersTx(container)

	filesHandler := grpc.NewFilesHandler(controller)

	files.RegisterFilesServer(svc.RPC(), filesHandler)

	return nil
}
