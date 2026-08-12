package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"log/slog"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/billing"
	"github.com/bd878/gallery/server/billing/internal/domain"
	"github.com/bd878/gallery/server/billing/pkg/events"
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
		domain.InvoicePayedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	slog.Debug("handle event", slog.String("name", event.EventName()), slog.Any("id", event.ID()), slog.Any("payload", event.Payload()))

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

	data, err := proto.Marshal(&billing.InvoicePayed{
		Id:       payload.ID,
		UserId:   payload.UserID,
		Cart:     payload.Cart,
		PayedAt:  payload.PayedAt,
	})
	if err != nil {
		return err
	}

	return h.stream.Publish(ctx, events.BillingChannel, am.NewEventMessage(event.ID(), events.InvoicePayedEvent, data, event.Metadata()))
}

func (h domainHandler[T]) onPremiumPayed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.PremiumPayed)

	data, err := proto.Marshal(&billing.PremiumPayed{
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

	return h.stream.Publish(ctx, events.BillingChannel, am.NewEventMessage(event.ID(), events.PremiumPayedEvent, data, event.Metadata()))
}