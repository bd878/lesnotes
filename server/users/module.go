package users

import (
	"context"
	"log/slog"

	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/ddd"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/users/internal/handler/stream"
	"github.com/bd878/gallery/server/users/config"
	users "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
	sessionsgateway "github.com/bd878/gallery/server/users/internal/gateway/sessions/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/service"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log).WithLang(httpmiddleware.Language)

	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS())

	outboxStore := pg.NewOutboxStore(svc.Pool(), "users_stream.outbox")

	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(
			am.NewEventStream(
				am.RawMessageStreamWithMiddleware(
					js,
					tm.NewOutboxStreamMiddleware(outboxStore),
				),
			),
		))

	ctrl := controller.New(cfg, sessionsGateway, dispatcher)

	stream.RegisterIntegrationEventHandlers(
		js,
		stream.NewIntegrationEventHandlers(ctrl),
	)

	handler := httphandler.New(ctrl, httphandler.Config{
		CookieDomain:    cfg.CookieDomain,
	})

	// TODO: middleware.Build(handler, ...middlewares)
	middleware.WithAuth(httpmiddleware.AuthBuilder(ctrl, sessionsGateway, users.PublicUserID))
	svc.ServeMux().Handle("/users/v1/me",     middleware.Build(handler.GetMe))
	svc.ServeMux().Handle("/users/v1/logout", middleware.Build(handler.Logout))
	svc.ServeMux().Handle("/users/v1/update", middleware.Build(handler.Update))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(ctrl, sessionsGateway, users.PublicUserID))
	svc.ServeMux().Handle("/users/v2/delete", middleware.Build(handler.DeleteJsonAPI))
	svc.ServeMux().Handle("/users/v2/me",     middleware.Build(handler.GetMe))
	svc.ServeMux().Handle("/users/v2/update", middleware.Build(handler.UpdateJsonAPI))

	middleware.NoAuth()
	svc.ServeMux().Handle("/users/v1/signup", middleware.Build(handler.Signup))
	svc.ServeMux().Handle("/users/v1/login",  middleware.Build(handler.Login))
	svc.ServeMux().Handle("/users/v1/auth",   middleware.Build(handler.Auth))
	svc.ServeMux().Handle("/liveness", middleware.Build(handler.Status))
	svc.ServeMux().Handle("/users/v2/signup", middleware.Build(handler.SignupJsonAPI))
	svc.ServeMux().Handle("/users/v2/auth",   middleware.Build(handler.AuthJsonAPI))
	svc.ServeMux().Handle("/users/v2/login",  middleware.Build(handler.LoginJsonAPI))

	return nil
}

func startOutboxProcessor(ctx context.Context, stream am.MessagePublisher[am.RawMessage], db pg.DB) {
	store := pg.NewOutboxStore(db, "users_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}