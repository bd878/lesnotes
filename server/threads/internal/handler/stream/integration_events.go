package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"log/slog"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/messages"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
)

type MessagesController interface {
	PrivateMessages(ctx context.Context, ids []int64, userID int64) error
	PublishMessages(ctx context.Context, ids []int64, userID int64) error
}

type ThreadsController interface {
	CreateThread(ctx context.Context, id, userID, parentID int64, name, description, title string, private bool) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
}

type integrationHandlers struct {
	messages   MessagesController
	threads    ThreadsController
}

var _ am.MessageHandler[am.EventMessage] = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController, threads ThreadsController) am.MessageHandler[am.EventMessage] {
	return integrationHandlers{
		messages:    messages,
		threads:     threads,
	}
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.EventMessage) error {
	slog.Debug("handle message", slog.String("name", msg.MessageName()))

	switch msg.MessageName() {
	case messageevents.MessagesPublishEvent:
		return h.handleMessagesPublish(ctx, msg)
	case messageevents.MessagesPrivateEvent:
		return h.handleMessagesPrivate(ctx, msg)
	case messageevents.MessageCreatedEvent:
		return h.handleMessageCreated(ctx, msg)
	case messageevents.MessageDeletedEvent:
		return h.handleMessageDeleted(ctx, msg)
	}

	return nil
}

// TODO: delete user threads when user deleted

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.CreateThread(ctx, m.Id, m.UserId, m.ThreadId, m.Name, "", "", m.Private)
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.DeleteThread(ctx, m.Id, m.UserId)
}

func (h integrationHandlers) handleMessagesPublish(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PublishMessages(ctx, m.GetIds(), m.GetUserId())
}

func (h integrationHandlers) handleMessagesPrivate(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PrivateMessages(ctx, m.GetIds(), m.GetUserId())
}
