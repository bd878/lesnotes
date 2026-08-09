package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/billing"
	"github.com/bd878/gallery/server/billing/internal/domain"
	"github.com/bd878/gallery/server/billing/pkg/model"
	"github.com/bd878/gallery/server/db/billing/pkg/loadbalance"
	"github.com/bd878/gallery/server/db/billing/pkg/machine"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/rpc"
)

type Config struct {
	RpcAddr string
}

type Controller struct {
	conf      Config
	client    billing.BillingClient
	conn      *grpc.ClientConn
	publisher ddd.EventPublisher[ddd.Event]
}

func New(conf Config, publisher ddd.EventPublisher[ddd.Event]) *Controller {
	controller := &Controller{conf: conf, publisher: publisher}

	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			controller.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := billing.NewBillingClient(conn)
	controller.conn = conn
	controller.client = client

	return controller
}

func (s *Controller) CreateInvoice(ctx context.Context, id string, userID int64, total int64, metadata []byte, cart *model.Cart) (err error) {
	slog.Debug("create invoice", slog.String("id", id), slog.Int64("user_id", userID), slog.Int64("total", total), slog.Any("metadata", metadata), slog.Any("cart", cart))

	cc, err := model.CartToProto(cart)
	if err != nil {
		return err
	}

	cartProto, err := proto.Marshal(cc)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&billing.AppendInvoiceCommand{
		Id:        id,
		UserId:    userID,
		Total:     total,
		Cart:      cartProto,
		Status:    "unpaid",
		Metadata:  metadata,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendInvoiceRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error) {
	slog.Debug("start payment", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("invoice_id", invoiceID), slog.String("currency", currency), slog.Int64("total", total), slog.Any("metadata", metadata))

	cmd, err := proto.Marshal(&billing.AppendPaymentCommand{
		Id:        id,
		UserId:    userID,
		InvoiceId: invoiceID,
		Status:    "pending",
		Currency:  currency,
		Total:     total,
		Metadata:  metadata,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendPaymentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("proceed payment", slog.Int64("id", id), slog.Int64("user_id", userID))

	payment, err := s.GetPayment(ctx, id, userID)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&billing.ProceedPaymentCommand{
		Id:        id,
		UserId:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	ii, err := s.GetInvoice(ctx, payment.InvoiceID, userID)
	if err != nil {
		return err
	}

	invoice, err := model.InvoiceToProto(ii)
	if err != nil {
		return err
	}

	slog.Debug("get invoice", slog.Int64("payment_id", id), slog.String("invoice_id", invoice.Id), slog.Int64("user_id", userID))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	events := make([]ddd.Event, 0, len(invoice.Cart.Items))

	invoiceEvent, err := domain.PayInvoice(invoice.Id, invoice.Cart, userID, updatedAt)
	if err != nil {
		return err
	}

	events = append(events, invoiceEvent)

	for _, cartItem := range invoice.Cart.Items {
		switch v := cartItem.Item.(type) {
		case *billing.CartItem_Premium:
			event, err := domain.PayPremium(invoice.Id, userID,
				v.Premium.ExpiresAt, updatedAt, v.Premium.Cost, v.Premium.Discount)
			if err != nil {
				return err
			}

			events = append(events, event)
		}
	}

	cmd1, err := proto.Marshal(&billing.PayInvoiceCommand{
		Id:        payment.InvoiceID,
		UserId:    userID,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.ProceedPaymentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})
	if err != nil {
		// TODO: rollback payment if failed
		return
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.PayInvoiceRequest),
		Cmd:      cmd1,
		Duration: "10s",
	})
	if err != nil {
		// TODO: rollback payment if failed
		return err
	}

	return s.publisher.Publish(ctx, events...)
}

func (s *Controller) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("cancel payment", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&billing.CancelPaymentCommand{
		Id:        id,
		UserId:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.CancelPaymentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("refund payment", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&billing.RefundPaymentCommand{
		Id:        id,
		UserId:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.RefundPaymentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) GetInvoice(ctx context.Context, id string, userID int64) (invoice *model.Invoice, err error) {
	slog.Debug("get invoice", slog.String("id", id), slog.Int64("user_id", userID))

	resp, err := s.client.GetInvoice(ctx, &billing.GetInvoiceRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	invoice, err = model.InvoiceFromProto(resp.Invoice)
	if err != nil {
		return nil, err
	}

	return
}

func (s *Controller) GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error) {
	slog.Debug("get payment", slog.Int64("id", id), slog.Int64("user_id", userID))

	resp, err := s.client.GetPayment(ctx, &billing.GetPaymentRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	payment = model.PaymentFromProto(resp.Payment)

	return
}
