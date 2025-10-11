package stream

import (
	"context"

	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/logger"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
)

type integrationHandlers struct {
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers() am.RawMessageHandler {
	return integrationHandlers{}
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
	logger.Debugw("handle message created", "name", msg.MessageName())
	return nil
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message deleted", "name", msg.MessageName())
	return nil
}