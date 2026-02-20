package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/api"
	filesevents "github.com/bd878/gallery/server/files/pkg/events"
)

type FilesController interface {
	DeleteFile(ctx context.Context, id, userID int64) error
}

type integrationHandlers struct {
	files       FilesController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(files FilesController) am.RawMessageHandler {
	return integrationHandlers{
		files: files,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(filesevents.FilesChannel, handlers)
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case filesevents.FileDeletedEvent:
		return h.handleFileDeleted(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleFileDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FileDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.DeleteFile(ctx, m.GetId(), m.GetUserId())
}
