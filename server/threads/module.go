package threads

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/rpc"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/threads/config"
	"github.com/bd878/gallery/server/threads/internal/handler/stream"
	"github.com/bd878/gallery/server/threads/internal/controller/service"
	"github.com/bd878/gallery/server/db/threads/pkg/loadbalance"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/threads/internal/handler/http"
	"github.com/bd878/gallery/server/internal/system"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	container := di.New()

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				loadbalance.Name,
				cfg.ThreadsServiceAddr,
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
	container.AddSingleton("threadsClient", func(c di.Container) (any, error) {
		client := threads.NewThreadsClient(c.Get("conn").(*grpc.ClientConn))
		return client, nil
	})
	container.AddSingleton("domainDispatcher", func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})
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
	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		outboxStore := pg.NewOutboxStore(tx, "threads_stream.outbox")
		return am.RawMessageStreamWithMiddleware(
			c.Get("js").(am.RawMessageStream),
			tm.NewOutboxStreamMiddleware(outboxStore),
		), nil
	})
	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		return am.NewEventStream(c.Get("txStream").(am.RawMessageStream)), nil
	})
	container.AddScoped("replyStream", func(c di.Container) (any, error) {
		return am.NewReplyStream(c.Get("txStream").(am.RawMessageStream)), nil
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		inboxStore := pg.NewInboxStore(tx, "threads_stream.inbox")
		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		return stream.NewDomainEventHandlers(c.Get("eventStream").(am.EventStream)), nil
	})

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	stream.RegisterDomainEventHandlersTx(container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event]))
	if err = stream.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, container)

	ctrl := service.NewThreadsControllerTx(
		container,
		service.New(container, container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])),
	)

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return stream.NewIntegrationEventHandlers(ctrl, ctrl), nil
	})
	container.AddScoped("commandHandlers", func(c di.Container) (any, error) {
		return stream.NewCommandHandlers(ctrl), nil
	})

	stream.RegisterIntegrationEventHandlersTx(container)

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/threads/v1/publish", middleware.Build(handler.PublishThread))
	svc.ServeMux().Handle("/threads/v1/private", middleware.Build(handler.PrivateThread))
	svc.ServeMux().Handle("/threads/v1/reorder", middleware.Build(handler.ReorderThread))
	svc.ServeMux().Handle("PUT /threads/v1/update", middleware.Build(handler.UpdateThread))

	middleware.NoAuth()
	svc.ServeMux().Handle("/liveness", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/threads/v2/read", middleware.Build(handler.ReadThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/list", middleware.Build(handler.ListThreadsJsonAPI))
	svc.ServeMux().Handle("/threads/v2/publish", middleware.Build(handler.PublishThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/private", middleware.Build(handler.PrivateThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/resolve", middleware.Build(handler.ResolveThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/create", middleware.Build(handler.CreateThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/update", middleware.Build(handler.UpdateThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/delete", middleware.Build(handler.DeleteThreadJsonAPI))
	svc.ServeMux().Handle("/threads/v2/reorder", middleware.Build(handler.ReorderThreadJsonAPI))

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	db := container.Get("db").(pg.DB)
	stream := container.Get("js").(am.MessagePublisher[am.RawMessage])

	store := pg.NewOutboxStore(db, "threads_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}
