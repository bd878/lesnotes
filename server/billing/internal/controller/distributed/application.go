package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/billing/pkg/model"
	"github.com/bd878/gallery/server/billing/internal/machine"
	"github.com/bd878/gallery/server/billing/internal/domain"
)

// TODO: can we return *api.Payment directly from repo w/o model conversion?
type PaymentsRepository interface {
	GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error)
}

type InvoicesRepository interface {
	GetInvoice(ctx context.Context, id string, userID int64) (invoice *api.Invoice, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus      Consensus
	log            *logger.Logger
	publisher      ddd.EventPublisher[ddd.Event]
	paymentsRepo   PaymentsRepository
	invoicesRepo   InvoicesRepository
}

func New(consensus Consensus, publisher ddd.EventPublisher[ddd.Event], paymentsRepo PaymentsRepository,
	invoicesRepo InvoicesRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:            log,
		publisher:      publisher,
		consensus:      consensus,
		paymentsRepo:   paymentsRepo,
		invoicesRepo:   invoicesRepo,
	}
}

func (m *Distributed) apply(ctx context.Context, reqType machine.RequestType, cmd []byte) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), 10*time.Second)
}

func (m *Distributed) CreateInvoice(ctx context.Context, id string, userID int64, total int64, metadata []byte, cart *api.Cart) (err error) {
	m.log.Debugw("invoice payment", "id", id, "user_id", userID, "total", total, "metadata", metadata, "cart", cart)

	cc, err := proto.Marshal(cart)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.AppendInvoiceCommand{
		Id:            id,
		UserId:        userID,
		Total:         total,
		Cart:          cc,
		Status:        "unpaid",
		Metadata:      metadata,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.AppendInvoiceRequest, cmd)
}

func (m *Distributed) StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error) {
	m.log.Debugw("start payment", "id", id, "user_id", userID, "invoice_id", invoiceID, "currency", currency, "total", total, "metadata", metadata)

	cmd, err := proto.Marshal(&machine.AppendPaymentCommand{
		Id:            id,
		UserId:        userID,
		InvoiceId:     invoiceID,
		Status:        "pending",
		Currency:      currency,
		Total:         total,
		Metadata:      metadata,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.AppendPaymentRequest, cmd)
}

func (m *Distributed) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("proceed payment", "id", id, "user_id", userID)

	payment, err := m.paymentsRepo.GetPayment(ctx, id, userID)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.ProceedPaymentCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	invoice, err := m.invoicesRepo.GetInvoice(ctx, payment.InvoiceID, userID)
	if err != nil {
		return err
	}

	m.log.Debugw("get invoice", "payment_id", id, "invoice_id", invoice.Id, "user_id", userID)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	events := make([]ddd.Event, 0)

	invoiceEvent, err := domain.PayInvoice(invoice.Id, invoice.Cart, userID, updatedAt)
	if err != nil {
		return err
	}

	events = append(events, invoiceEvent)

	for _, cartItem := range invoice.Cart.Items {
		switch v := cartItem.Item.(type) {
		case *api.CartItem_Premium:
			event, err := domain.PayPremium(invoice.Id, userID,
				v.Premium.ExpiresAt, updatedAt, v.Premium.Cost, v.Premium.Discount)
			if err != nil {
				return err
			}

			events = append(events, event)
		}
	}

	cmd1, err := proto.Marshal(&machine.PayInvoiceCommand{
		Id:            payment.InvoiceID,
		UserId:        userID,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.ProceedPaymentRequest, cmd)
	if err != nil {
		// TODO: rollback payment if failed
		return
	}

	err = m.apply(ctx, machine.PayInvoiceRequest, cmd1)
	if err != nil {
		// TODO: rollback payment if failed
		return err
	}

	return m.publisher.Publish(context.TODO(), events...)
}

func (m *Distributed) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("cancel payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.CancelPaymentCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.CancelPaymentRequest, cmd)
}

func (m *Distributed) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("refund payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.RefundPaymentCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.RefundPaymentRequest, cmd)
}

func (m *Distributed) GetInvoice(ctx context.Context, id string, userID int64) (invoice *api.Invoice, err error) {
	m.log.Debugw("get invoice", "id", id, "user_id", userID)
	return m.invoicesRepo.GetInvoice(ctx, id, userID)
}

func (m *Distributed) GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error) {
	m.log.Debugw("get payment", "id", id, "user_id", userID)
	return m.paymentsRepo.GetPayment(ctx, id, userID)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}