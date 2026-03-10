package domain

import (
	"errors"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	InvoicePayedEvent = "billing.InvoicePayed"
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
