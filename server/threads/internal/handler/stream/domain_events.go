package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/ddd"
	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/threads/internal/domain"
	"github.com/bd878/gallery/server/threads/pkg/events"
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
		domain.ThreadCreatedEvent,
		domain.ThreadDeletedEvent,
		domain.ThreadPublishEvent,
		domain.ThreadPrivateEvent,
		domain.ThreadParentChangedEvent,
		domain.ThreadUpdatedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ThreadCreatedEvent:
		return h.onThreadCreated(ctx, event)
	case domain.ThreadDeletedEvent:
		return h.onThreadDeleted(ctx, event)
	case domain.ThreadUpdatedEvent:
		return h.onThreadUpdated(ctx, event)
	case domain.ThreadPrivateEvent:
		return h.onThreadPrivated(ctx, event)
	case domain.ThreadPublishEvent:
		return h.onThreadPublished(ctx, event)
	case domain.ThreadParentChangedEvent:
		return h.onThreadParentChanged(ctx, event)
	}
	return nil
}

func (h domainHandler[T]) onThreadCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadCreated)
	data, err := proto.Marshal(&api.ThreadCreated{
		Id:          payload.ID,
		UserId:      payload.UserID,
		ParentId:    payload.ParentID,
		Name:        payload.Name,
		Description: payload.Description,
		Private:     payload.Private,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadCreatedEvent, data))
}

func (h domainHandler[T]) onThreadDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadDeleted)
	data, err := proto.Marshal(&api.ThreadDeleted{
		Id:          payload.ID,
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadDeletedEvent, data))	
}

func (h domainHandler[T]) onThreadPrivated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadPrivated)
	data, err := proto.Marshal(&api.ThreadPrivated{
		Id:          payload.ID,
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadPrivatedEvent, data))	
}

func (h domainHandler[T]) onThreadUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadUpdated)
	data, err := proto.Marshal(&api.ThreadUpdated{
		Id:          payload.ID,
		UserId:      payload.UserID,
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadUpdatedEvent, data))	
}

func (h domainHandler[T]) onThreadPublished(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadPublished)
	data, err := proto.Marshal(&api.ThreadPublished{
		Id:          payload.ID,
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadPublishedEvent, data))	
}

func (h domainHandler[T]) onThreadParentChanged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadParentChanged)
	data, err := proto.Marshal(&api.ThreadParentChanged{
		Id:          payload.ID,
		UserId:      payload.UserID,
		ParentId:    payload.ParentID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.MessagesChannel, am.NewRawMessage(event.ID(), events.ThreadParentChangedEvent, data))		
}
