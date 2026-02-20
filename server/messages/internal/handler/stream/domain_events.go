package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/messages/pkg/events"
)

type domainHandler[T ddd.Event] struct {
	publisher am.MessagePublisher[am.Message]
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[am.Message]) *domainHandler[ddd.Event] {
	return &domainHandler[ddd.Event]{publisher}
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
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	logger.Debugw("handle event", "name", event.EventName(), "id", event.ID(), "payload", event.Payload())

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
	}
	return nil
}

func (h domainHandler[T]) onMessageCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageCreated)
	data, err := proto.Marshal(&api.MessageCreated{
		Id:       payload.ID,
		UserId:   payload.UserID,
		Text:     payload.Text,
		Title:    payload.Title,
		Name:     payload.Name,
		Private:  payload.Private,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.MessageCreatedEvent, data))
}

func (h domainHandler[T]) onMessageDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageDeleted)
	data, err := proto.Marshal(&api.MessageDeleted{
		Id:     payload.ID,
		UserId: payload.UserID, 
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.MessageDeletedEvent, data))
}

func (h domainHandler[T]) onMessageUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageUpdated)
	data, err := proto.Marshal(&api.MessageUpdated{
		Id:       payload.ID,
		UserId:   payload.UserID,
		Text:     payload.Text,
		Title:    payload.Title,
		Name:     payload.Name,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.MessageUpdatedEvent, data))
}

func (h domainHandler[T]) onMessagesPrivate(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessagesPrivated)
	data, err := proto.Marshal(&api.MessagesPrivated{
		Ids:      payload.IDs,
		UserId:   payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.MessagesPrivateEvent, data))
}

func (h domainHandler[T]) onMessagesPublish(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessagesPublished)
	data, err := proto.Marshal(&api.MessagesPublished{
		Ids:       payload.IDs,
		UserId:    payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.MessagesPublishEvent, data))
}

func (h domainHandler[T]) onTranslationCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationCreated)
	data, err := proto.Marshal(&api.TranslationCreated{
		MessageId:    payload.MessageID,
		UserId:       payload.UserID,
		Lang:         payload.Lang,
		Text:         payload.Text,
		Title:        payload.Title,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.TranslationsChannel, am.NewRawMessage(event.ID(), events.TranslationCreatedEvent, data))
}

func (h domainHandler[T]) onTranslationDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationDeleted)
	data, err := proto.Marshal(&api.TranslationDeleted{
		MessageId:     payload.MessageID,
		Lang:          payload.Lang,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.TranslationsChannel, am.NewRawMessage(event.ID(), events.TranslationDeletedEvent, data))
}

func (h domainHandler[T]) onTranslationUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TranslationUpdated)
	data, err := proto.Marshal(&api.TranslationUpdated{
		MessageId:        payload.MessageID,
		Lang:             payload.Lang,
		Text:             payload.Text,
		Title:            payload.Title,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.TranslationsChannel, am.NewRawMessage(event.ID(), events.TranslationUpdatedEvent, data))
}
