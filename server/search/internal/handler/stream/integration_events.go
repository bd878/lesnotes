package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/api"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
	threadsevents "github.com/bd878/gallery/server/threads/pkg/events"
	filesevents "github.com/bd878/gallery/server/files/pkg/events"
)

type MessagesController interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) error
	UpdateMessage(ctx context.Context, id, userID int64, name, title, text string) error
	PublishMessages(ctx context.Context, ids []int64, userID int64) error
	PrivateMessages(ctx context.Context, ids []int64, userID int64) error
	DeleteMessage(ctx context.Context, id, userID int64) error
}

type ThreadsController interface {
	SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool) error
	DeleteThread(ctx context.Context, id, userID int64) error
	UpdateThread(ctx context.Context, id, userID int64, name, description string) error
	ChangeThreadParent(ctx context.Context, id, userID, parentID int64) error
	PrivateThread(ctx context.Context, id, userID int64) error
	PublishThread(ctx context.Context, id, userID int64) error
}

type FilesController interface {
	SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64) error
	PublishFile(ctx context.Context, id, userID int64) error
	PrivateFile(ctx context.Context, id, userID int64) error
	DeleteFile(ctx context.Context, id, userID int64) error
}

type integrationHandlers struct {
	messages MessagesController
	threads  ThreadsController
	files    FilesController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController, threads ThreadsController, files FilesController) am.RawMessageHandler {
	return integrationHandlers{
		messages: messages,
		threads:  threads,
		files:    files,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(messageevents.MessagesChannel, handlers)
	if err != nil {
		return
	}

	err = subscriber.Subscribe(threadsevents.MessagesChannel, handlers)
	if err != nil {
		return
	}

	err = subscriber.Subscribe(filesevents.FilesChannel, handlers)
	if err != nil {
		return
	}

	return
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

	case threadsevents.ThreadCreatedEvent:
		return h.handleThreadCreated(ctx, msg)
	case threadsevents.ThreadDeletedEvent:
		return h.handleThreadDeleted(ctx, msg)
	case threadsevents.ThreadParentChangedEvent:
		return h.handleThreadParentChanged(ctx, msg)
	case threadsevents.ThreadUpdatedEvent:
		return h.handleThreadUpdated(ctx, msg)
	case threadsevents.ThreadPublishedEvent:
		return h.handleThreadPublished(ctx, msg)
	case threadsevents.ThreadPrivatedEvent:
		return h.handleThreadPrivated(ctx, msg)

	case filesevents.FileUploadedEvent:
		return h.handleFileUploaded(ctx, msg)
	case filesevents.FileDeletedEvent:
		return h.handleFileDeleted(ctx, msg)
	case filesevents.FilePublishedEvent:
		return h.handleFilePublished(ctx, msg)
	case filesevents.FilePrivatedEvent:
		return h.handleFilePrivated(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.SaveMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetTitle(), m.GetText(), m.GetPrivate())
}

func (h integrationHandlers) handleMessageDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.DeleteMessage(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleMessageUpdated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.UpdateMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetTitle(), m.GetText())
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

func (h integrationHandlers) handleThreadCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.SaveThread(ctx, m.GetId(), m.GetUserId(), m.GetParentId(), m.GetName(), m.GetDescription(), m.GetPrivate())
}

func (h integrationHandlers) handleThreadDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.DeleteThread(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleThreadUpdated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.UpdateThread(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetDescription())
}

func (h integrationHandlers) handleThreadParentChanged(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.ChangeThreadParent(ctx, m.GetId(), m.GetUserId(), m.GetParentId())
}

func (h integrationHandlers) handleThreadPrivated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.PrivateThread(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleThreadPublished(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.PublishThread(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleFileUploaded(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FileUploaded{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.SaveFile(ctx, m.GetId(), m.GetUserId(), m.GetName(), m.GetDescription(), m.GetMime(), m.GetPrivate(), m.GetSize())
}

func (h integrationHandlers) handleFileDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FileDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.DeleteFile(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleFilePublished(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FilePublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PublishFile(ctx, m.GetId(), m.GetUserId())
}

func (h integrationHandlers) handleFilePrivated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FilePrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PrivateFile(ctx, m.GetId(), m.GetUserId())
}
