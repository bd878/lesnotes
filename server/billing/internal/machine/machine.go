package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/billing/pkg/model"
)

type PaymentsRepository interface {
	SavePayment(ctx context.Context, id, userID int64, invoiceID string, currency, status string, total int64, metadata []byte) (err error)
	ProceedPayment(ctx context.Context, id, userID int64) (err error)
	CancelPayment(ctx context.Context, id, userID int64) (err error)
	RefundPayment(ctx context.Context, id, userID int64) (err error)
	GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type InvoicesRepository interface {
	SaveInvoice(ctx context.Context, id string, userID int64, currency, status string, total int64, metadata []byte) (err error)
	CancelInvoice(ctx context.Context, id string, userID int64) (err error)
	PayInvoice(ctx context.Context, id string, userID int64) (err error)
	GetInvoice(ctx context.Context, id string, userID int64) (invoice *model.Invoice, err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log          *logger.Logger
	paymentsRepo PaymentsRepository
	invoicesRepo InvoicesRepository
}

func New(paymentsRepo PaymentsRepository, invoicesRepo InvoicesRepository, log *logger.Logger) *Machine {
	return &Machine{
		log:          log,
		paymentsRepo: paymentsRepo,
		invoicesRepo: invoicesRepo,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendInvoiceRequest:
		return f.applyAppendInvoice(buf[1:])
	case AppendPaymentRequest:
		return f.applyAppendPayment(buf[1:])
	case ProceedPaymentRequest:
		return f.applyProceedPayment(buf[1:])
	case PayInvoiceRequest:
		return f.applyPayInvoice(buf[1:])
	case CancelPaymentRequest:
		return f.applyCancelPayment(buf[1:])
	case CancelInvoiceRequest:
		return f.applyCancelInvoice(buf[1:])
	case RefundPaymentRequest:
		return f.applyRefundPayment(buf[1:])
	default:
		f.log.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppendInvoice(raw []byte) interface{} {
	var cmd AppendInvoiceCommand
	proto.Unmarshal(raw, &cmd)

	return f.invoicesRepo.SaveInvoice(context.Background(), cmd.Id, cmd.UserId, cmd.Currency, cmd.Status, cmd.Total, cmd.Metadata)
}

func (f *Machine) applyAppendPayment(raw []byte) interface{} {
	var cmd AppendPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.SavePayment(context.Background(), cmd.Id, cmd.UserId, cmd.InvoiceId,
		cmd.Currency, cmd.Status, cmd.Total, cmd.Metadata)
}

func (f *Machine) applyProceedPayment(raw []byte) interface{} {
	var cmd ProceedPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.ProceedPayment(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPayInvoice(raw []byte) interface{} {
	var cmd PayInvoiceCommand
	proto.Unmarshal(raw, &cmd)

	return f.invoicesRepo.PayInvoice(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyCancelPayment(raw []byte) interface{} {
	var cmd CancelPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.CancelPayment(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyCancelInvoice(raw []byte) interface{} {
	var cmd CancelInvoiceCommand
	proto.Unmarshal(raw, &cmd)

	return f.invoicesRepo.CancelInvoice(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyRefundPayment(raw []byte) interface{} {
	var cmd RefundPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.RefundPayment(context.Background(), cmd.Id, cmd.UserId)
}
