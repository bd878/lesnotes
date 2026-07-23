package billing

import (
	"context"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/billing/internal/handler/stream"

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

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher, stream.NewDomainEventHandlers(
		am.NewMessagePublisher(
			jetstream.NewStream(svc.Config().NatsStream, svc.JS(), svc.Logger()),
		),
	))

	ctrl := controller.New(controller.Config{RpcAddr: cfg.BillingServiceAddr}, dispatcher)

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("GET   /billing/v1/invoices", middleware.Build(handler.GetInvoice))
	svc.ServeMux().Handle("GET   /billing/v1/payments", middleware.Build(handler.GetPayment))

	middleware.NoAuth()
	svc.ServeMux().Handle("GET   /liveness", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("POST  /billing/v2/invoices", middleware.Build(handler.CreateInvoiceJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/payments", middleware.Build(handler.StartPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/cancel",   middleware.Build(handler.CancelPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/proceed",  middleware.Build(handler.ProceedPaymentJsonAPI))
	svc.ServeMux().Handle("POST  /billing/v2/refund",   middleware.Build(handler.RefundPaymentJsonAPI))

	return nil
}
