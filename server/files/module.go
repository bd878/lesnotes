package files

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/amotel"
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
	messagesRepo := postgres.NewMessagesRepository(svc.Pool(), "files.messages")

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(am.NewMessagePublisher(
			nats.NewStream(svc.Nats()),
			amotel.OtelMessageContextInjector(),
		)))

	controller := application.New(dispatcher, filesRepo, messagesRepo, svc.Logger())

	stream.RegisterIntegrationEventHandlers(
		am.NewMessageSubscriber(
			nats.NewStream(svc.Nats()),
			amotel.OtelMessageContextExtractor(),
		),
		stream.NewIntegrationEventHandlers(controller),
	)

	filesHandler := grpc.NewFilesHandler(controller)

	api.RegisterFilesServer(svc.RPC(), filesHandler)

	return nil
}