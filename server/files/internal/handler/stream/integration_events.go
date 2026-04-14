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
	SaveMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) error
	DeleteMessageFiles(ctx context.Context, id, userID int64) error
	UpdateMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) error
	PublishMessageFiles(ctx context.Context, userID int64, messageIDs []int64) error
	PrivateMessageFiles(ctx context.Context, userID int64, messageIDs []int64) error
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

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.SaveMessageFiles(ctx, m.GetId(), m.GetUserId(), m.GetFileIds())
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.DeleteMessageFiles(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleMessageUpdated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.UpdateMessageFiles(ctx, m.GetId(), m.GetUserId(), m.GetFileIds())
}

func (h integrationHandlers) handleMessagesPublished(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PublishMessageFiles(ctx, m.GetUserId(), m.GetIds())
}

func (h integrationHandlers) handleMessagesPrivated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PrivateMessageFiles(ctx, m.GetUserId(), m.GetIds())
}
