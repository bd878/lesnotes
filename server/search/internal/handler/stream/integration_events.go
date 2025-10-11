package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/api"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
)

type Controller interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string) error
	DeleteMessage(ctx context.Context, id, userID int64) error
}

type integrationHandlers struct {
	controller Controller
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(controller Controller) am.RawMessageHandler {
	return integrationHandlers{controller}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	return subscriber.Subscribe(messageevents.MessagesChannel, handlers)
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case messageevents.MessageCreatedEvent:
		return h.handleMessageCreated(ctx, msg)
	case messageevents.MessageDeletedEvent:
		return h.handleMessageDeleted(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.SaveMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetTitle(), m.GetText())
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.DeleteMessage(ctx, m.GetId(), m.GetUserId())
}