package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/users/internal/domain"
	"github.com/bd878/gallery/server/users/pkg/events"
)

type domainHandler[T ddd.Event] struct {
	publisher am.MessagePublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher) *domainHandler[ddd.Event] {
	return &domainHandler[ddd.Event]{publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handler ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handler,
		domain.UserDeletedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.UserDeletedEvent:
		return h.onUserDeleted(ctx, event)
	}
	return nil
}

func (h domainHandler[T]) onUserDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.UserDeleted)
	data, err := proto.Marshal(&users.UserDeleted{
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.UsersChannel, am.NewRawMessage(event.ID(), events.UserDeletedEvent, data, event.Metadata(), events.UsersChannel))	
}
