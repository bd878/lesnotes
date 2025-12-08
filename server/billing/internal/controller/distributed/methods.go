package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	billing "github.com/bd878/gallery/server/billing/pkg/model"
)

func (m *Distributed) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	timeout := 10*time.Second
	/* fsm.Apply() */
	future := m.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res = future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}

	return
}

func (m *Distributed) CreateInvoice(ctx context.Context, id string, userID int64, currency string, total int64, metadata []byte) (err error) {
	logger.Debugw("invoice payment", "id", id, "user_id", userID, "currency", currency, "total", total, "metadata", metadata)

	cmd, err := proto.Marshal(&AppendInvoiceCommand{
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

	_, err = m.apply(ctx, AppendInvoiceRequest, cmd)

	return
}

func (m *Distributed) StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error) {
	logger.Debugw("start payment", "id", id, "user_id", userID, "invoice_id", invoiceID, "currency", currency, "total", total, "metadata", metadata)

	cmd, err := proto.Marshal(&AppendPaymentCommand{
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

	_, err = m.apply(ctx, AppendPaymentRequest, cmd)

	return
}

func (m *Distributed) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("proceed payment", "id", id, "user_id", userID)

	payment, err := m.paymentsRepo.GetPayment(ctx, id, userID)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&ProceedPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, ProceedPaymentRequest, cmd)

	cmd1, err := proto.Marshal(&PayInvoiceCommand{
		Id:     payment.InvoiceID,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PayInvoiceRequest, cmd1)

	return
}

func (m *Distributed) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("cancel payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&CancelPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, CancelPaymentRequest, cmd)

	return
}

func (m *Distributed) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("refund payment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&RefundPaymentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, RefundPaymentRequest, cmd)

	return
}

func (m *Distributed) GetInvoice(ctx context.Context, id string, userID int64) (invoice *billing.Invoice, err error) {
	logger.Debugw("get invoice", "id", id, "user_id", userID)
	return m.invoicesRepo.GetInvoice(ctx, id, userID)
}

func (m *Distributed) GetPayment(ctx context.Context, id, userID int64) (payment *billing.Payment, err error) {
	logger.Debugw("get payment", "id", id, "user_id", userID)
	return m.paymentsRepo.GetPayment(ctx, id, userID)
}
