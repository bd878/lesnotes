package tm

import (
	"context"
	"errors"
	"github.com/bd878/gallery/server/internal/am"
)

type OutboxStore interface {
	Save(ctx context.Context, msg am.Message) error
	FindUnpublished(ctx context.Context, limit int) ([]am.Message, error)
	MarkPublished(ctx context.Context, ids ...string) error
}

type outbox struct {
	am.MessagePublisher
	store OutboxStore
}

var _ am.MessagePublisher = (*outbox)(nil)

func NewOutboxStreamMiddleware(store OutboxStore) am.MessagePublisherMiddleware {
	o := outbox{store: store}

	return func(publisher am.MessagePublisher) am.MessagePublisher {
		o.MessagePublisher = publisher

		return o
	}
}

func (o outbox) Publish(ctx context.Context, topicName string, msg am.Message) error {
	err := o.store.Save(ctx, msg)

	var errDupe ErrDuplicateMessage
	if errors.Is(err, errDupe) {
		return nil
	}

	return err
}