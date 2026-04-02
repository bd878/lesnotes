package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/api"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
)

type MessagesController interface {
	PrivateMessages(ctx context.Context, ids []int64, userID int64) error
	PublishMessages(ctx context.Context, ids []int64, userID int64) error
}

type integrationHandlers struct {
	log        *logger.Logger
	messages   MessagesController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController, log *logger.Logger) am.RawMessageHandler {
	return integrationHandlers{
		log:         log,
		messages:    messages,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(messageevents.MessagesChannel, handlers)
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	h.log.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case messageevents.MessagesPublishEvent:
		return h.handleMessagesPublish(ctx, msg)
	case messageevents.MessagesPrivateEvent:
		return h.handleMessagesPrivate(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessagesPublish(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PublishMessages(ctx, m.GetIds(), m.GetUserId())
}

func (h integrationHandlers) handleMessagesPrivate(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PrivateMessages(ctx, m.GetIds(), m.GetUserId())
}
