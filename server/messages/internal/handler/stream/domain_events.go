package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/ddd"
	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/api"
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
	subscriber.Subscribe(handler, domain.MessageCreatedEvent, domain.MessageDeletedEvent)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.MessageCreatedEvent:
		return h.onMessageCreated(ctx, event)
	case domain.MessageDeletedEvent:
		return h.onMessageDeleted(ctx, event)
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
