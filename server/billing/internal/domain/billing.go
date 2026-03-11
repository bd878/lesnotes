package domain

import (
	"errors"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	InvoicePayedEvent = "billing.InvoicePayed"
	PremiumPayedEvent = "billing.PremiumPayed"
)

var (
	ErrIDRequired = errors.New("id is empty")
)

type InvoicePayed struct {
	ID        string
	UserID    int64
	Cart      *api.Cart
	PayedAt   string
}

func (InvoicePayed) Key() string { return InvoicePayedEvent }

func PayInvoice(id string, cart *api.Cart, userID int64, payedAt string) (ddd.Event, error) {
	if id == "" {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(InvoicePayedEvent, &InvoicePayed{
		ID:         id,
		Cart:       cart,
		UserID:     userID,
		PayedAt:    payedAt,
	}), nil
}

type PremiumPayed struct {
	InvoiceID  string
	UserID     int64
	ExpiresAt  string
	Cost       int64
	Discount   int64
	CreatedAt  string
}

func (PremiumPayed) Key() string { return PremiumPayedEvent }

func PayPremium(invoiceID string, userID int64, expiresAt, createdAt string, cost, discount int64) (ddd.Event, error) {
	if invoiceID == "" {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(PremiumPayedEvent, &PremiumPayed{
		InvoiceID:  invoiceID,
		UserID:     userID,
		ExpiresAt:  expiresAt,
		CreatedAt:  createdAt,
		Cost:       cost,
		Discount:   discount,
	}), nil
}