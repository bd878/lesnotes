package stream

import (
	"context"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/ddd"
)

func RegisterDomainEventHandlersTx(subscriber ddd.EventSubscriber[ddd.Event]) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		domainHandlers := di.Get(ctx, "domainEventHandlers").(ddd.EventHandler[ddd.Event])

		return domainHandlers.HandleEvent(ctx, event)
	})

	RegisterDomainEventHandlers(subscriber, handlers)
}
