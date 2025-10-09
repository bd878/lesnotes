package stream

import (
	"context"

	"github.com/bd878/gallery/server/ddd"
	"github.com/bd878/gallery/server/am"
	"github.com/bd878/gallery/server/messages/"
)

type integrationHandlers[T ddd.Event] struct {
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers() ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	if err = subscriber.Subscribe(); err != nil {
		return
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case 
	}
}

func (h integrationHandlers[T]) handleMessageCreated(ctx context.Context, event T) error {
	return nil
}

func (h integrationHandlers[T]) handleMessageDeleted(ctx context.Context, event T) error {
	return nil
}