package billing

import (
	"log/slog"
	"context"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/billing/internal/handler/stream"
	pg "github.com/bd878/gallery/server/internal/postgres"

	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/billing/internal/handler/http"
	controller "github.com/bd878/gallery/server/billing/internal/controller/service"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	outboxStore := pg.NewOutboxStore(svc.Pool(), "billing_stream.outbox")

	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS())

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher, stream.NewDomainEventHandlers(
		am.NewEventStream(
			am.RawMessageStreamWithMiddleware(
				js,
				tm.NewOutboxStreamMiddleware(outboxStore),
			),
		),
	))

	startOutboxProcessor(ctx, js, svc.Pool())

	ctrl := controller.New(controller.Config{RpcAddr: cfg.BillingServiceAddr}, dispatcher)

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

func startOutboxProcessor(ctx context.Context, stream am.MessagePublisher[am.RawMessage], db pg.DB) {
	store := pg.NewOutboxStore(db, "billing_stream.outbox")
	outboxProcessor := tm.NewOutboxProcessor(stream, store)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			slog.Error("failed to start outbox processor", slog.String("err", err.Error()))
		}
	}()
}
