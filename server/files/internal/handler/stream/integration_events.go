package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	messagesevents "github.com/bd878/gallery/server/messages/pkg/events"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
)

type FilesController interface {

}

type integrationHandlers struct {
	files FilesController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(files FilesController) am.RawMessageHandler {
	return integrationHandlers{
		files: files,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(messagesevents.MessagesChannel, handlers)
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case messagesevents.MessagesPublishEvent:
		return h.handleMessagesPublished(ctx, msg)
	case messagesevents.MessagesPrivateEvent:
		return h.handleMessagesPrivated(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessagesPublished(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	// TODO: implement
	return nil
}

func (h integrationHandlers) handleMessagesPrivated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	// TODO: implement
	return nil
}
