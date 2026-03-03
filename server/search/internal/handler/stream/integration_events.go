package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/api"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
	threadsevents "github.com/bd878/gallery/server/threads/pkg/events"
	filesevents "github.com/bd878/gallery/server/files/pkg/events"
)

type MessagesController interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool, createdAt, updatedAt string) error
	UpdateMessage(ctx context.Context, id, userID int64, name, title, text *string, updatedAt string) error
	PublishMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) error
	PrivateMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) error
	DeleteMessage(ctx context.Context, id, userID int64) error
}

type ThreadsController interface {
	SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool, createdAt, updatedAt string) error
	DeleteThread(ctx context.Context, id, userID int64) error
	UpdateThread(ctx context.Context, id, userID int64, name, description *string, updatedAt string) error
	ChangeThreadParent(ctx context.Context, id, userID, parentID int64) error
	PrivateThread(ctx context.Context, id, userID int64, updatedAt string) error
	PublishThread(ctx context.Context, id, userID int64, updatedAt string) error
}

type FilesController interface {
	SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64, createdAt, updatedAt string) error
	PublishFile(ctx context.Context, id, userID int64, updatedAt string) error
	PrivateFile(ctx context.Context, id, userID int64, updatedAt string) error
	DeleteFile(ctx context.Context, id, userID int64) error
}

type TranslationsController interface {
	SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string, createdAt, updatedAt string) error
	DeleteTranslation(ctx context.Context, messageID int64, lang string) error
	UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string, updatedAt string) error
}

type integrationHandlers struct {
	log            *logger.Logger
	messages       MessagesController
	threads        ThreadsController
	files          FilesController
	translations   TranslationsController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController, threads ThreadsController,
	files FilesController, translations TranslationsController, log *logger.Logger) am.RawMessageHandler {
	return integrationHandlers{
		log:            log,
		messages:       messages,
		translations:   translations,
		threads:        threads,
		files:          files,
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

	err = subscriber.Subscribe(messageevents.TranslationsChannel, handlers)
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	h.log.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

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

	case messageevents.TranslationCreatedEvent:
		return h.handleTranslationCreated(ctx, msg)
	case messageevents.TranslationDeletedEvent:
		return h.handleTranslationDeleted(ctx, msg)
	case messageevents.TranslationUpdatedEvent:
		return h.handleTranslationUpdated(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleMessageCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessageCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.SaveMessage(ctx, m.GetId(), m.GetUserId(), m.GetName(),
		m.GetTitle(), m.GetText(), m.GetPrivate(), m.GetCreatedAt(), m.GetUpdatedAt())
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

	return h.messages.UpdateMessage(ctx, m.GetId(), m.GetUserId(), m.Name, m.Title, m.Text, m.GetUpdatedAt())
}

func (h integrationHandlers) handleMessagesPublish(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PublishMessages(ctx, m.GetIds(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleMessagesPrivate(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.MessagesPrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.PrivateMessages(ctx, m.GetIds(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleThreadCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.SaveThread(ctx, m.GetId(), m.GetUserId(), m.GetParentId(),
		m.GetName(), m.GetDescription(), m.GetPrivate(), m.GetCreatedAt(), m.GetUpdatedAt())
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

	return h.threads.UpdateThread(ctx, m.GetId(), m.GetUserId(), m.Name, m.Description, m.GetUpdatedAt())
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

	return h.threads.PrivateThread(ctx, m.GetId(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleThreadPublished(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.ThreadPublished{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.threads.PublishThread(ctx, m.GetId(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleFileUploaded(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FileUploaded{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.SaveFile(ctx, m.GetId(), m.GetUserId(), m.GetName(),
		m.GetDescription(), m.GetMime(), m.GetPrivate(), m.GetSize(), m.GetCreatedAt(), m.GetUpdatedAt())
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

	return h.files.PublishFile(ctx, m.GetId(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleFilePrivated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.FilePrivated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.files.PrivateFile(ctx, m.GetId(), m.GetUserId(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleTranslationCreated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.TranslationCreated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.translations.SaveTranslation(ctx, m.GetUserId(), m.GetMessageId(),
		m.GetLang(), m.GetTitle(), m.GetText(), m.GetCreatedAt(), m.GetUpdatedAt())
}

func (h integrationHandlers) handleTranslationDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.TranslationDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.translations.DeleteTranslation(ctx, m.GetMessageId(), m.GetLang())
}

func (h integrationHandlers) handleTranslationUpdated(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.TranslationUpdated{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.translations.UpdateTranslation(ctx, m.GetMessageId(), m.GetLang(), m.Title, m.Text, m.GetUpdatedAt())
}
