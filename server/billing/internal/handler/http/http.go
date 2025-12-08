package http

import (
	"io"
	"context"
	"net/http"

	billingmodel "github.com/bd878/gallery/server/billing/pkg/model"
)

type Controller interface {
	CreateInvoice(ctx context.Context, id string, userID int64, currency string, total int64, metadata []byte) (err error)
	StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error)
	ProceedPayment(ctx context.Context, id, userID int64) (err error)
	CancelPayment(ctx context.Context, id, userID int64) (err error)
	RefundPayment(ctx context.Context, id, userID int64) (err error)
	GetInvoice(ctx context.Context, id string, userID int64) (invoice *billingmodel.Invoice, err error)
	GetPayment(ctx context.Context, id, userID int64) (payment *billingmodel.Payment, err error)
}

type Handler struct {
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
