package files

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/files/config"
	"github.com/bd878/gallery/server/files/internal/repository/postgres"
	"github.com/bd878/gallery/server/files/internal/controller/application"
	"github.com/bd878/gallery/server/files/internal/handler/grpc"
	"github.com/bd878/gallery/server/files/internal/handler/stream"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	filesRepo := postgres.NewFilesRepository(svc.Pool(), "files.files")

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(nats.NewStream(svc.Nats())))

	controller := application.New(dispatcher, filesRepo, svc.Logger())

	stream.RegisterIntegrationEventHandlers(nats.NewStream(svc.Nats()),
		stream.NewIntegrationEventHandlers(controller))

	filesHandler := grpc.NewFilesHandler(controller)

	api.RegisterFilesServer(svc.RPC(), filesHandler)

	return nil
}