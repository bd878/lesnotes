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
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) error
	UpdateMessage(ctx context.Context, id, userID int64, name, title, text string, private int32) error
	PublishMessages(ctx context.Context, ids []int64, userID int64) error
	PrivateMessages(ctx context.Context, ids []int64, userID int64) error
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
	case messageevents.MessageUpdatedEvent:
		return h.handleMessageUpdated(ctx, msg)
	case messageevents.MessagesPublishEvent:
		return h.handleMessagesPublish(ctx, msg)
	case messageevents.MessagesPrivateEvent:
		return h.handleMessagesPrivate(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.SaveMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetTitle(), m.GetText(), m.GetPrivate())
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.DeleteMessage(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleMessageUpdated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.UpdateMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetTitle(), m.GetText(), m.GetPrivate())
}

func (h integrationHandlers) handleMessagesPublish(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.PublishMessages(ctx, m.GetIds(), m.GetUserId())
}

func (h integrationHandlers) handleMessagesPrivate(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.controller.PrivateMessages(ctx, m.GetIds(), m.GetUserId())
}
