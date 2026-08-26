package saga

import (
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/threads"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/sec"
	"github.com/bd878/gallery/server/messages/internal/domain"
)

type sagaHandler[T ddd.Event] struct {
	orchestrator sec.Orchestrator
}

func NewEventHandlers(saga sec.Orchestrator) *sagaHandler[ddd.Event] {
	return &sagaHandler[ddd.Event]{orchestrator: saga}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event]) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		sagaEventHandlers := di.Get(ctx, "sagaEventHandlers").(ddd.EventHandler[ddd.Event])

		return sagaEventHandlers.HandleEvent(ctx, event)
	})

	subscriber.Subscribe(handlers,
		domain.MessageCreatedEvent,
		domain.MessageDeletedEvent,
	)
}

func (h sagaHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	switch event.EventName() {
	case domain.MessageCreatedEvent:
		return h.onMessageCreated(ctx, event)
	}
	return nil
}

func (h sagaHandler[T]) onMessageCreated(ctx context.Context, event ddd.Event) error {
	slog.Debug("handle domain command event",
		slog.String("name", event.EventName()),
		slog.String("id", event.ID()),
	)

	payload := event.Payload().(*domain.MessageCreated)
	data, err := proto.Marshal(&threads.CreateThread{
		ThreadId: payload.ID,
		UserId: payload.UserID,
		ParentId: payload.ThreadID,
		Private: &payload.Private,
		Name: payload.Name,
	})
	if err != nil {
		return err
	}

	return h.orchestrator.Start(ctx, event.ID(), data)
}
