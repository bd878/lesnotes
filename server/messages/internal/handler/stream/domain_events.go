package stream

import (
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/internal/am"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/messages/pkg/events"
)

type domainHandler[T ddd.Event] struct {
	stream am.EventStream
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(stream am.EventStream) *domainHandler[ddd.Event] {
	return &domainHandler[ddd.Event]{stream}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handler ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handler,
		domain.MessageCreatedEvent,
		domain.MessageDeletedEvent,
		domain.MessagesPrivateEvent,
		domain.MessagesPublishEvent,
		domain.MessageUpdatedEvent,

		domain.TranslationCreatedEvent,
		domain.TranslationDeletedEvent,
		domain.TranslationUpdatedEvent,

		domain.CommentCreatedEvent,
		domain.CommentDeletedEvent,
		domain.CommentUpdatedEvent,
		domain.MessageCommentsDeletedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	slog.Debug("handle event", slog.String("name", event.EventName()), slog.Any("id", event.ID()), slog.Any("payload", event.Payload()))

	switch event.EventName() {
	case domain.MessageCreatedEvent:
		return h.onMessageCreated(ctx, event)
	case domain.MessageDeletedEvent:
		return h.onMessageDeleted(ctx, event)
	case domain.MessageUpdatedEvent:
		return h.onMessageUpdated(ctx, event)
	case domain.MessagesPrivateEvent:
		return h.onMessagesPrivate(ctx, event)
	case domain.MessagesPublishEvent:
		return h.onMessagesPublish(ctx, event)

	case domain.TranslationCreatedEvent:
		return h.onTranslationCreated(ctx, event)
	case domain.TranslationUpdatedEvent:
		return h.onTranslationUpdated(ctx, event)
	case domain.TranslationDeletedEvent:
		return h.onTranslationDeleted(ctx, event)

	case domain.CommentCreatedEvent:
		return h.onCommentCreated(ctx, event)
	case domain.CommentUpdatedEvent:
		return h.onCommentUpdated(ctx, event)
	case domain.CommentDeletedEvent:
		return h.onCommentDeleted(ctx, event)
	case domain.MessageCommentsDeletedEvent:
		return h.onMessageCommentsDeleted(ctx, event)
	}
	return nil
}

func (h domainHandler[T]) onMessageCreated(ctx context.Context, event ddd.Event) error {
	// TODO: add serde, create message in am EventPublisher : am/event_messages.go
	payload := event.Payload().(*domain.MessageCreated)
	data, err := proto.Marshal(&messages.MessageCreated{
		Id:        payload.ID,
		UserId:    payload.UserID,
		FileIds:   payload.FileIDs,
		Text:      payload.Text,
		Title:     payload.Title,
		Name:      payload.Name,
		Private:   payload.Private,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.MessagesChannel, am.NewEventMessage(event.ID(), events.MessageCreatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onMessageDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageDeleted)
	data, err := proto.Marshal(&messages.MessageDeleted{
		Id:     payload.ID,
		UserId: payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.MessagesChannel, am.NewEventMessage(event.ID(), events.MessageDeletedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onMessageUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageUpdated)
	data, err := proto.Marshal(&messages.MessageUpdated{
		Id:        payload.ID,
		UserId:    payload.UserID,
		FileIds:   payload.FileIDs,
		Text:      payload.Text,
		Title:     payload.Title,
		Name:      payload.Name,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.MessagesChannel, am.NewEventMessage(event.ID(), events.MessageUpdatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onMessagesPrivate(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessagesPrivated)
	data, err := proto.Marshal(&messages.MessagesPrivated{
		Ids:       payload.IDs,
		UserId:    payload.UserID,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.MessagesChannel, am.NewEventMessage(event.ID(), events.MessagesPrivateEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onMessagesPublish(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessagesPublished)
	data, err := proto.Marshal(&messages.MessagesPublished{
		Ids:       payload.IDs,
		UserId:    payload.UserID,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.MessagesChannel, am.NewEventMessage(event.ID(), events.MessagesPublishEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onTranslationCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationCreated)
	data, err := proto.Marshal(&translations.TranslationCreated{
		Id:        payload.MessageID,
		UserId:    payload.UserID,
		Lang:      payload.Lang,
		Text:      payload.Text,
		Title:     payload.Title,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.TranslationsChannel, am.NewEventMessage(event.ID(), events.TranslationCreatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onTranslationDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationDeleted)
	data, err := proto.Marshal(&translations.TranslationDeleted{
		Id:        payload.MessageID,
		Lang:      payload.Lang,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.TranslationsChannel, am.NewEventMessage(event.ID(), events.TranslationDeletedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onTranslationUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationUpdated)
	data, err := proto.Marshal(&translations.TranslationUpdated{
		Id:        payload.MessageID,
		Lang:      payload.Lang,
		Text:      payload.Text,
		Title:     payload.Title,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.TranslationsChannel, am.NewEventMessage(event.ID(), events.TranslationUpdatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onCommentCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CommentCreated)
	data, err := proto.Marshal(&comments.CommentCreated{
		MessageId: payload.MessageID,
		Id:        payload.ID,
		UserId:    payload.UserID,
		Text:      payload.Text,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.CommentsChannel, am.NewEventMessage(event.ID(), events.CommentCreatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onCommentUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CommentUpdated)
	data, err := proto.Marshal(&comments.CommentUpdated{
		Id:        payload.ID,
		UserId:    payload.UserID,
		Text:      payload.Text,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.CommentsChannel, am.NewEventMessage(event.ID(), events.CommentUpdatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onCommentDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CommentDeleted)
	data, err := proto.Marshal(&comments.CommentDeleted{
		Id:     payload.ID,
		UserId: payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.CommentsChannel, am.NewEventMessage(event.ID(), events.CommentDeletedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onMessageCommentsDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageCommentsDeleted)
	data, err := proto.Marshal(&comments.MessageCommentsDeleted{
		MessageId: payload.MessageID,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.CommentsChannel, am.NewEventMessage(event.ID(), events.MessageCommentsDeletedEvent, data, event.Metadata()))
}
