package events

import (
	"context"

	"github.com/bd878/gallery/server/ddd"
	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/internal/domain"
)

type domainHandler[T ddd.Event] struct {
	publisher am.MessagePublisher[am.RawMessage]
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[am.RawMessage]) *domainHandler[ddd.Event] {
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
	logger.Debugw("message created event", "name", event.EventName())
	return nil
}

func (h domainHandler[T]) onMessageDeleted(ctx context.Context, event ddd.Event) error {
	logger.Debugw("message deleted event", "name", event.EventName())
	return nil
}
