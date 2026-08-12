package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/threads/internal/domain"
	"github.com/bd878/gallery/server/threads/pkg/events"
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
	data, err := proto.Marshal(&threads.ThreadCreated{
		Id:          payload.ID,
		UserId:      payload.UserID,
		ParentId:    payload.ParentID,
		Name:        payload.Name,
		Description: payload.Description,
		Title:       payload.Title,
		Private:     payload.Private,
		CreatedAt:   payload.CreatedAt,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadCreatedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onThreadDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadDeleted)
	data, err := proto.Marshal(&threads.ThreadDeleted{
		Id:          payload.ID,
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadDeletedEvent, data, event.Metadata()))	
}

func (h domainHandler[T]) onThreadPrivated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadPrivated)
	data, err := proto.Marshal(&threads.ThreadPrivated{
		Id:          payload.ID,
		UserId:      payload.UserID,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadPrivatedEvent, data, event.Metadata()))	
}

func (h domainHandler[T]) onThreadUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadUpdated)
	data, err := proto.Marshal(&threads.ThreadUpdated{
		Id:          payload.ID,
		UserId:      payload.UserID,
		Name:        payload.Name,
		Description: payload.Description,
		Title:       payload.Title,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadUpdatedEvent, data, event.Metadata()))	
}

func (h domainHandler[T]) onThreadPublished(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadPublished)
	data, err := proto.Marshal(&threads.ThreadPublished{
		Id:          payload.ID,
		UserId:      payload.UserID,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadPublishedEvent, data, event.Metadata()))	
}

func (h domainHandler[T]) onThreadParentChanged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ThreadParentChanged)
	data, err := proto.Marshal(&threads.ThreadParentChanged{
		Id:          payload.ID,
		UserId:      payload.UserID,
		ParentId:    payload.ParentID,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.ThreadsChannel, am.NewEventMessage(event.ID(), events.ThreadParentChangedEvent, data, event.Metadata()))		
}
