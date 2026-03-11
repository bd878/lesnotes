package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type PaymentsRepository interface {
	SavePayment(ctx context.Context, id, userID int64, invoiceID string, currency, status string, total int64, metadata []byte, createdAt, updatedAt string) (err error)
	ProceedPayment(ctx context.Context, id, userID int64, updatedAt string) (err error)
	CancelPayment(ctx context.Context, id, userID int64, updatedAt string) (err error)
	RefundPayment(ctx context.Context, id, userID int64, updatedAt string) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type InvoicesRepository interface {
	SaveInvoice(ctx context.Context, id string, userID int64, status string, total int64, metadata []byte, cart []byte, createdAt, updatedAt string) (err error)
	CancelInvoice(ctx context.Context, id string, userID int64, updatedAt string) (err error)
	PayInvoice(ctx context.Context, id string, userID int64, updatedAt string) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type Dumper interface {
	Open(ctx context.Context) (ch chan *api.BillingSnapshot, err error)
	Restore(ctx context.Context, snapshot *api.BillingSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log             *logger.Logger
	dumper          Dumper
	paymentsRepo    PaymentsRepository
	invoicesRepo    InvoicesRepository
}

func New(paymentsRepo PaymentsRepository, invoicesRepo InvoicesRepository, dumper Dumper, log *logger.Logger) *Machine {
	return &Machine{
		log:          log,
		dumper:       dumper,
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

	return f.invoicesRepo.SaveInvoice(context.TODO(), cmd.Id, cmd.UserId, cmd.Status, cmd.Total,
		cmd.Metadata, cmd.Cart, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyAppendPayment(raw []byte) interface{} {
	var cmd AppendPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.SavePayment(context.TODO(), cmd.Id, cmd.UserId, cmd.InvoiceId,
		cmd.Currency, cmd.Status, cmd.Total, cmd.Metadata, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyProceedPayment(raw []byte) interface{} {
	var cmd ProceedPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.ProceedPayment(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyPayInvoice(raw []byte) interface{} {
	var cmd PayInvoiceCommand
	proto.Unmarshal(raw, &cmd)

	return f.invoicesRepo.PayInvoice(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyCancelPayment(raw []byte) interface{} {
	var cmd CancelPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.CancelPayment(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyCancelInvoice(raw []byte) interface{} {
	var cmd CancelInvoiceCommand
	proto.Unmarshal(raw, &cmd)

	return f.invoicesRepo.CancelInvoice(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyRefundPayment(raw []byte) interface{} {
	var cmd RefundPaymentCommand
	proto.Unmarshal(raw, &cmd)

	return f.paymentsRepo.RefundPayment(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}
