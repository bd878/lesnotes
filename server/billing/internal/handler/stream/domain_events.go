package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/billing/internal/domain"
	"github.com/bd878/gallery/server/billing/pkg/events"
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
		domain.InvoicePayedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	logger.Debugw("handle event", "name", event.EventName(), "id", event.ID(), "payload", event.Payload())

	switch event.EventName() {
	case domain.InvoicePayedEvent:
		return h.onInvoicePayed(ctx, event)
	case domain.PremiumPayedEvent:
		return h.onPremiumPayed(ctx, event)
	}

	return nil
}

func (h domainHandler[T]) onInvoicePayed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoicePayed)

	data, err := proto.Marshal(&api.InvoicePayed{
		Id:       payload.ID,
		UserId:   payload.UserID,
		Cart:     payload.Cart,
		PayedAt:  payload.PayedAt,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.BillingChannel, am.NewRawMessage(event.ID(), events.InvoicePayedEvent, data))
}

func (h domainHandler[T]) onPremiumPayed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.PremiumPayed)

	data, err := proto.Marshal(&api.PremiumPayed{
		InvoiceId:   payload.InvoiceID,
		ExpiresAt:   payload.ExpiresAt,
		CreatedAt:   payload.CreatedAt,
		UserId:      payload.UserID,
		Cost:        payload.Cost,
		Discount:    payload.Discount,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.BillingChannel, am.NewRawMessage(event.ID(), events.PremiumPayedEvent, data))
}