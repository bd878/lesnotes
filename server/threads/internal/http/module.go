package http

import (
	"context"

	"github.com/bd878/gallery/server/threads/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/threads/internal/handler/stream"

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

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(
			am.NewMessagePublisher(
				nats.NewStream(svc.Nats()),
			),
		))

	ctrl := controller.New(cfg, dispatcher)

	handler := httphandler.New(ctrl)

	stream.RegisterIntegrationEventHandlers(
		am.NewMessageSubscriber(
			nats.NewStream(svc.Nats()),
		),
		stream.NewIntegrationEventHandlers(ctrl, ctrl, svc.Logger()),
	)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/threads/v1/publish", middleware.Build(handler.PublishThread))
	svc.ServeMux().Handle("/threads/v1/private", middleware.Build(handler.PrivateThread))
	svc.ServeMux().Handle("/threads/v1/reorder", middleware.Build(handler.ReorderThread))
	svc.ServeMux().Handle("PUT /threads/v1/update", middleware.Build(handler.UpdateThread))

	middleware.NoAuth()
	svc.ServeMux().Handle("/threads/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
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
