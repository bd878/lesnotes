package search

import (
	"context"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/jetstream"
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
	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)

	ctrl := controller.New(controller.Config{RpcAddr: cfg.SearchServiceAddr})

	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS(), svc.Logger())
	stream.RegisterIntegrationEventHandlers(
		am.NewMessageSubscriber(
			js,
		),
		stream.NewIntegrationEventHandlers(ctrl, ctrl, ctrl, ctrl, svc.Logger()))

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/search/v1/messages", middleware.Build(handler.SearchMessages))
	svc.ServeMux().Handle("/search/v1/files", middleware.Build(handler.SearchFiles))
	// TODO: search translations

	middleware.NoAuth()
	svc.ServeMux().Handle("/liveness", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/search/v2/messages", middleware.Build(handler.SearchMessagesJsonAPI))
	svc.ServeMux().Handle("/search/v2/files", middleware.Build(handler.SearchFilesJsonAPI))

	return nil
}
