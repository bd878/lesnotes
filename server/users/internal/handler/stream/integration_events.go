package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/api"
	billing "github.com/bd878/gallery/server/billing/pkg/events"
)

type UsersController interface {
	MakePremium(ctx context.Context, userID int64, invoiceID, createdAt, expiresAt string) error
}

type integrationHandlers struct {
	users     UsersController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(users UsersController) am.RawMessageHandler {
	return integrationHandlers{
		users:   users,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(billing.BillingChannel, handlers)
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case billing.PremiumPayedEvent:
		return h.handlePremiumPayed(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handlePremiumPayed(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.PremiumPayed{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.users.MakePremium(ctx, m.GetUserId(), m.GetInvoiceId(), m.GetCreatedAt(), m.GetExpiresAt())
}