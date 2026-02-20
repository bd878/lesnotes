package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/billing/pkg/model"
	"github.com/bd878/gallery/server/billing/internal/machine"
)

type PaymentsRepository interface {
	GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error)
}

type InvoicesRepository interface {
	GetInvoice(ctx context.Context, id string, userID int64) (invoice *model.Invoice, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus      Consensus
	log            *logger.Logger
	paymentsRepo   PaymentsRepository
	invoicesRepo   InvoicesRepository
}

func New(consensus Consensus, paymentsRepo PaymentsRepository, invoicesRepo InvoicesRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:            log,
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

func (m *Distributed) CreateInvoice(ctx context.Context, id string, userID int64, currency string, total int64, metadata []byte) (err error) {
	m.log.Debugw("invoice payment", "id", id, "user_id", userID, "currency", currency, "total", total, "metadata", metadata)

	cmd, err := proto.Marshal(&machine.AppendInvoiceCommand{
		Id:        id,
		UserId:    userID,
		Currency:  currency,
		Total:     total,
		Status:    "unpaid",
		Metadata:  metadata,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendInvoiceRequest, cmd)

	return
}

func (m *Distributed) StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error) {
	m.log.Debugw("start payment", "id", id, "user_id", userID, "invoice_id", invoiceID, "currency", currency, "total", total, "metadata", metadata)

	cmd, err := proto.Marshal(&machine.AppendPaymentCommand{
		Id:        id,
		UserId:    userID,
		InvoiceId: invoiceID,
		Status:    "pending",
		Currency:  currency,
		Total:     total,
		Metadata:  metadata,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendPaymentRequest, cmd)

	return
}

func (m *Distributed) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("proceed payment", "id", id, "user_id", userID)

	payment, err := m.paymentsRepo.GetPayment(ctx, id, userID)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.ProceedPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.ProceedPaymentRequest, cmd)
	if err != nil {
		return
	}

	cmd1, err := proto.Marshal(&machine.PayInvoiceCommand{
		Id:     payment.InvoiceID,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PayInvoiceRequest, cmd1)

	return
}

func (m *Distributed) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("cancel payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.CancelPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.CancelPaymentRequest, cmd)

	return
}

func (m *Distributed) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("refund payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.RefundPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.RefundPaymentRequest, cmd)

	return
}

func (m *Distributed) GetInvoice(ctx context.Context, id string, userID int64) (invoice *model.Invoice, err error) {
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