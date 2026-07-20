package files

import (
	"context"

	"github.com/bd878/gallery/server/api/files"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/jetstream"
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
	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS(), svc.Logger())
	stream.RegisterDomainEventHandlers(dispatcher,
		stream.NewDomainEventHandlers(am.NewMessagePublisher(
			js,
		)))

	controller := application.New(dispatcher, filesRepo, messagesRepo, svc.Logger())

	stream.RegisterIntegrationEventHandlers(
		am.NewMessageSubscriber(
			js,
		),
		stream.NewIntegrationEventHandlers(controller),
	)

	filesHandler := grpc.NewFilesHandler(controller)

	files.RegisterFilesServer(svc.RPC(), filesHandler)

	return nil
}