package billing

import (
	"fmt"
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/di"
	pg "github.com/bd878/gallery/server/internal/postgres"

	"github.com/bd878/gallery/server/api/billing"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/db/billing/pkg/loadbalance"
	"github.com/bd878/gallery/server/billing/internal/handler/stream"
	"github.com/bd878/gallery/server/billing/internal/controller/service"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/billing/internal/handler/http"
	"github.com/bd878/gallery/server/internal/system"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	container := di.New()

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				loadbalance.Name,
				cfg.BillingServiceAddr,
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
	container.AddSingleton("billingClient", func(c di.Container) (any, error) {
		client := billing.NewBillingClient(c.Get("conn").(*grpc.ClientConn))
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

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		tx := c.Get("tx").(pgx.Tx)
		outboxStore := pg.NewOutboxStore(tx, "billing_stream.outbox")

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

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlersTx(dispatcher)

	startOutboxProcessor(ctx, container)

	ctrl := service.NewBillingControllerTx(
		container,
		service.New(container, dispatcher),
	)

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("GET   /billing/v1/invoices", middleware.Build(handler.GetInvoice))
	svc.ServeMux().Handle("GET   /billing/v1/payments", middleware.Build(handler.GetPayment))

	middleware.NoAuth()
	svc.ServeMux().Handle("GET   /liveness", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("POST  /billing/v2/invoices", middleware.Build(handler.CreateInvoiceJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/payments", middleware.Build(handler.StartPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/cancel",   middleware.Build(handler.CancelPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/proceed",  middleware.Build(handler.ProceedPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/refund",   middleware.Build(handler.RefundPaymentJsonAPI))

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	db := container.Get("db").(pg.DB)
	stream := container.Get("js").(am.MessagePublisher[am.RawMessage])

	store := pg.NewOutboxStore(db, "billing_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}
