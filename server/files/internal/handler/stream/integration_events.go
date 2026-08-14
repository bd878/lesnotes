package stream

import (
	"log/slog"
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/messages"
	messagesevents "github.com/bd878/gallery/server/messages/pkg/events"
)

type FilesController interface {
	SaveMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) error
	DeleteMessageFiles(ctx context.Context, id, userID int64) error
	UpdateMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) error
	PublishMessageFiles(ctx context.Context, userID int64, messageIDs []int64) error
	PrivateMessageFiles(ctx context.Context, userID int64, messageIDs []int64) error
}

type integrationHandlers struct {
	files FilesController
}

var _ am.MessageHandler[am.EventMessage] = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(files FilesController) am.MessageHandler[am.EventMessage] {
	return integrationHandlers{
		files: files,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) error {
	return subscriber.Subscribe(messagesevents.MessagesChannel, handlers, am.GroupName("files-messages"))
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.EventMessage) error {
	slog.Debug("handle message", slog.String("name", msg.MessageName()))

	switch msg.MessageName() {
	case messagesevents.MessageCreatedEvent:
		return h.handleMessageCreated(ctx, msg)
	case messagesevents.MessageUpdatedEvent:
		return h.handleMessageUpdated(ctx, msg)
	case messagesevents.MessageDeletedEvent:
		return h.handleMessageDeleted(ctx, msg)
	case messagesevents.MessagesPublishEvent:
		return h.handleMessagesPublished(ctx, msg)
	case messagesevents.MessagesPrivateEvent:
		return h.handleMessagesPrivated(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.SaveMessageFiles(ctx, m.GetId(), m.GetUserId(), m.GetFileIds())
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.DeleteMessageFiles(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleMessageUpdated(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessageUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.UpdateMessageFiles(ctx, m.GetId(), m.GetUserId(), m.GetFileIds())
}

func (h integrationHandlers) handleMessagesPublished(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PublishMessageFiles(ctx, m.GetUserId(), m.GetIds())
}

func (h integrationHandlers) handleMessagesPrivated(ctx context.Context, msg am.EventMessage) error {
	m := &messages.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PrivateMessageFiles(ctx, m.GetUserId(), m.GetIds())
}
