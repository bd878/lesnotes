package search

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/di"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/search/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/search/internal/handler/stream"

	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/search/internal/handler/http"
	controller "github.com/bd878/gallery/server/search/internal/controller/service"
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
		inboxStore := pg.NewInboxStore(tx, "search_stream.inbox")
		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	ctrl := controller.New(controller.Config{RpcAddr: cfg.SearchServiceAddr})

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return stream.NewIntegrationEventHandlers(ctrl, ctrl, ctrl, ctrl), nil
	})

	stream.RegisterIntegrationEventHandlersTx(container)

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/search/v1/messages", middleware.Build(handler.SearchMessages))
	svc.ServeMux().Handle("/search/v1/files", middleware.Build(handler.SearchFiles))
	// TODO: search translations

	middleware.NoAuth()
	svc.ServeMux().Handle("/liveness", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/search/v2/messages", middleware.Build(handler.SearchMessagesJsonAPI))
	svc.ServeMux().Handle("/search/v2/files", middleware.Build(handler.SearchFilesJsonAPI))

	return nil
}
