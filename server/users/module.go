package users

import (
	"fmt"
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/users/internal/handler/stream"
	"github.com/bd878/gallery/server/users/config"
	"github.com/bd878/gallery/server/db/users/pkg/loadbalance"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/users/internal/handler/http"
	sessionsgateway "github.com/bd878/gallery/server/users/internal/gateway/sessions/grpc"
	controller "github.com/bd878/gallery/server/users/internal/controller/service"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	container := di.New()

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				loadbalance.Name,
				cfg.UsersServiceAddr,
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
	container.AddSingleton("usersClient", func(c di.Container) (any, error) {
		client := users.NewUsersClient(c.Get("conn").(*grpc.ClientConn))
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

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		inboxStore := pg.NewInboxStore(tx, "users_inbox.inbox")
		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		outboxStore := pg.NewOutboxStore(tx, "users_stream.outbox")

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

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log).WithLang(httpmiddleware.Language)

	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlersTx(dispatcher)

	startOutboxProcessor(ctx, container)

	ctrl := controller.NewUsersControllerTx(
		container,
		controller.New(container, sessionsGateway, dispatcher),
	)

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return stream.NewIntegrationEventHandlers(ctrl), nil
	})

	stream.RegisterIntegrationEventHandlersTx(container)

	handler := httphandler.New(ctrl, httphandler.Config{
		CookieDomain: cfg.CookieDomain,
	})

	middleware.WithAuth(httpmiddleware.AuthBuilder(ctrl, sessionsGateway, usersmodel.PublicUserID))
	svc.ServeMux().Handle("/users/v1/me",      middleware.Build(handler.GetMe))
	svc.ServeMux().Handle("/users/v1/logout",  middleware.Build(handler.Logout))
	svc.ServeMux().Handle("/users/v1/update",  middleware.Build(handler.Update))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(ctrl, sessionsGateway, usersmodel.PublicUserID))
	svc.ServeMux().Handle("/users/v2/delete", middleware.Build(handler.DeleteJsonAPI))
	svc.ServeMux().Handle("/users/v2/me",      middleware.Build(handler.GetMe))
	svc.ServeMux().Handle("/users/v2/update",  middleware.Build(handler.UpdateJsonAPI))

	middleware.NoAuth()
	svc.ServeMux().Handle("/users/v1/signup", middleware.Build(handler.Signup))
	svc.ServeMux().Handle("/users/v1/login",  middleware.Build(handler.Login))
	svc.ServeMux().Handle("/users/v1/auth",   middleware.Build(handler.Auth))
	svc.ServeMux().Handle("/liveness",         middleware.Build(handler.Status))
	svc.ServeMux().Handle("/users/v2/signup",  middleware.Build(handler.SignupJsonAPI))
	svc.ServeMux().Handle("/users/v2/auth",    middleware.Build(handler.AuthJsonAPI))
	svc.ServeMux().Handle("/users/v2/login",   middleware.Build(handler.LoginJsonAPI))

	svc.ServeMux().Handle("/metrics", promhttp.Handler())

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	db := container.Get("db").(pg.DB)
	stream := container.Get("js").(am.MessagePublisher[am.RawMessage])

	store := pg.NewOutboxStore(db, "users_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}
