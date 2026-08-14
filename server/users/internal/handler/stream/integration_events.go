package stream

import (
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/billing"
	billingevents "github.com/bd878/gallery/server/billing/pkg/events"
)

type UsersController interface {
	MakePremium(ctx context.Context, userID int64, invoiceID, createdAt, expiresAt string) error
}

type integrationHandlers struct {
	users     UsersController
}

var _ am.MessageHandler[am.EventMessage] = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(users UsersController) am.MessageHandler[am.EventMessage] {
	return integrationHandlers{
		users:   users,
	}
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.EventMessage) error {
	slog.Debug("handle message", slog.String("name", msg.MessageName()))

	switch msg.MessageName() {
	case billingevents.PremiumPayedEvent:
		return h.handlePremiumPayed(ctx, msg)
	}

	return nil
}

// TODO: event ddd.Event
func (h integrationHandlers) handlePremiumPayed(ctx context.Context, msg am.EventMessage) error {
	m := &billing.PremiumPayed{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.users.MakePremium(ctx, m.GetUserId(), m.GetInvoiceId(), m.GetCreatedAt(), m.GetExpiresAt())
}