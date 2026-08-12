package threads

import (
	"log/slog"
	"context"

	"github.com/bd878/gallery/server/threads/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/threads/internal/handler/stream"
	pg "github.com/bd878/gallery/server/internal/postgres"

	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/threads/internal/handler/http"
	controller "github.com/bd878/gallery/server/threads/internal/controller/service"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	outboxStore := pg.NewOutboxStore(svc.Pool(), "threads_stream.outbox")

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS())
	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(
			am.NewEventStream(
				am.RawMessageStreamWithMiddleware(
					js,
					tm.NewOutboxStreamMiddleware(outboxStore),
				),
			),
		))

	startOutboxProcessor(ctx, js, svc.Pool())

	ctrl := controller.New(cfg, dispatcher)

	handler := httphandler.New(ctrl)

	stream.RegisterIntegrationEventHandlers(
		js,
		stream.NewIntegrationEventHandlers(ctrl, ctrl),
	)

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

func startOutboxProcessor(ctx context.Context, stream am.MessagePublisher[am.RawMessage], db pg.DB) {
	store := pg.NewOutboxStore(db, "threads_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}