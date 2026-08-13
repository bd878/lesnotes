package messages

import (
	"fmt"
	"log/slog"
	"context"

	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/rpc"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/db/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/internal/handler/stream"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	controller "github.com/bd878/gallery/server/messages/internal/controller/service"
	threadsgateway "github.com/bd878/gallery/server/messages/internal/gateway/threads/grpc"
	filesgateway "github.com/bd878/gallery/server/messages/internal/gateway/files/grpc"
	httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	container := di.New()

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				loadbalance.Name,
				cfg.MessagesServiceAddr,
			),
		)
	})
	container.AddSingleton("messagesClient", func(c di.Container) (any, error) {
		client := messages.NewMessagesClient(c.Get("conn").(*grpc.ClientConn))
		return client, nil
	})
	container.AddSingleton("translationsClient", func(c di.Container) (any, error) {
		client := translations.NewTranslationsClient(c.Get("conn").(*grpc.ClientConn))
		return client, nil
	})
	container.AddSingleton("commentsClient", func(c di.Container) (any, error) {
		client := comments.NewCommentsClient(c.Get("conn").(*grpc.ClientConn))
		return client, nil
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

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	threadsGateway := threadsgateway.New(cfg.ThreadsServiceAddr)
	filesGateway := filesgateway.New(cfg.FilesServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		outboxStore := pg.NewOutboxStore(tx, "messages_stream.outbox")

		js := c.Get("js").(am.RawMessageStream)

		return stream.NewDomainEventHandlers(
			am.NewEventStream(
				am.RawMessageStreamWithMiddleware(
					js,
					tm.NewOutboxStreamMiddleware(outboxStore),
				),
			),
		), nil
	})

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlersTx(dispatcher)

	startOutboxProcessor(ctx, container)

	messagesSaved := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_saved_count",
		},
	)

	messagesController := controller.NewMessagesControllerTx(
		container,
		controller.NewMessagesController(container, dispatcher, filesGateway, threadsGateway, messagesSaved),
	)
	translationsController := controller.NewTranslationsControllerTx(
		container,
		controller.NewTranslationsController(container, dispatcher),
	)
	commentsController := controller.NewCommentsControllerTx(
		container,
		controller.NewCommentsController(container, dispatcher),
	)

	stream.RegisterIntegrationEventHandlers(
		container.Get("js").(am.RawMessageStream),
		stream.NewIntegrationEventHandlers(messagesController),
	)

	handler := httphandler.New(messagesController, translationsController, commentsController)

	svc.ServeMux().Handle("/messages/v1/send", middleware.Build(handler.SendMessage))
	svc.ServeMux().Handle("/messages/v1/read_path", middleware.Build(handler.ReadPath))
	svc.ServeMux().Handle("/messages/v1/read", middleware.Build(handler.ReadMessages))
	svc.ServeMux().Handle("/messages/v1/update", middleware.Build(handler.UpdateMessage))
	svc.ServeMux().Handle("/messages/v1/publish", middleware.Build(handler.PublishMessages))
	svc.ServeMux().Handle("/messages/v1/private", middleware.Build(handler.PrivateMessages))
	svc.ServeMux().Handle("/messages/v1/delete", middleware.Build(handler.DeleteMessages))

	middleware.NoAuth()
	svc.ServeMux().Handle("GET /liveness", middleware.Build(handler.GetStatus))
	svc.ServeMux().Handle("GET /metrics", promhttp.Handler())

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/messages/v2/send", middleware.Build(handler.SendMessageJsonAPI))
	svc.ServeMux().Handle("/messages/v2/read", middleware.Build(handler.ReadMessagesJsonAPI))
	// TODO: /threads/v2/read
	svc.ServeMux().Handle("/messages/v2/read_path", middleware.Build(handler.ReadPathJsonAPI))
	svc.ServeMux().Handle("/messages/v2/read_tree", middleware.Build(handler.ReadTreeJsonAPI))
	svc.ServeMux().Handle("/messages/v2/publish", middleware.Build(handler.PublishMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/private", middleware.Build(handler.PrivateMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/delete", middleware.Build(handler.DeleteMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/update", middleware.Build(handler.UpdateMessageJsonAPI))

	svc.ServeMux().Handle("/translations/v2/send", middleware.Build(handler.SendTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/update", middleware.Build(handler.UpdateTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/delete", middleware.Build(handler.DeleteTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/read", middleware.Build(handler.ReadTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/list", middleware.Build(handler.ListTranslationsJsonAPI))

	svc.ServeMux().Handle("POST /comments/v2/send", middleware.Build(handler.SendCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/update", middleware.Build(handler.UpdateCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/read", middleware.Build(handler.ReadCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/delete", middleware.Build(handler.DeleteCommentsJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/list", middleware.Build(handler.ListCommentsJsonAPI))

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	db := container.Get("db").(pg.DB)
	stream := container.Get("js").(am.MessagePublisher[am.RawMessage])

	store := pg.NewOutboxStore(db, "messages_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}