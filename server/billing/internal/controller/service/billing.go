package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/billing/pkg/loadbalance"
	billingmodel "github.com/bd878/gallery/server/billing/pkg/model"
)

type Config struct {
	RpcAddr  string
}

type Controller struct {
	conf         Config
	client       api.BillingClient
	conn         *grpc.ClientConn
}

func New(conf Config) *Controller {
	controller := &Controller{conf: conf}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Controller) setupConnection() (err error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := api.NewBillingClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
		return true
	}
	return false
}

func (s *Controller) CreateInvoice(ctx context.Context, id string, userID int64, currency string, total int64, metadata []byte) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create invoice", "id", id, "user_id", userID, "currency", currency, "total", total, "metadata", metadata)

	_, err = s.client.CreateInvoice(ctx, &api.CreateInvoiceRequest{
		Id:        id,
		UserId:    userID,
		Currency:  currency,
		Total:     total,
		Metadata:  metadata,
	})

	return
}

func (s *Controller) StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("start payment", "id", id, "user_id", userID, "invoice_id", invoiceID, "currency", currency, "total", total, "metadata", metadata)

	_, err = s.client.StartPayment(ctx, &api.StartPaymentRequest{
		Id:        id,
		UserId:    userID,
		InvoiceId: invoiceID,
		Currency:  currency,
		Total:     total,
		Metadata:  metadata,
	})

	return
}

func (s *Controller) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("proceed payment", "id", id, "user_id", userID)

	_, err = s.client.ProceedPayment(ctx, &api.ProceedPaymentRequest{
		Id:        id,
		UserId:    userID,
	})

	return
}

func (s *Controller) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("cancel payment", "id", id, "user_id", userID)

	_, err = s.client.CancelPayment(ctx, &api.CancelPaymentRequest{
		Id:        id,
		UserId:    userID,
	})

	return
}

func (s *Controller) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("refund payment", "id", id, "user_id", userID)

	_, err = s.client.RefundPayment(ctx, &api.RefundPaymentRequest{
		Id:        id,
		UserId:    userID,
	})

	return
}

func (s *Controller) GetInvoice(ctx context.Context, id string, userID int64) (invoice *billingmodel.Invoice, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("get invoice", "id", id, "user_id", userID)

	resp, err := s.client.GetInvoice(ctx, &api.GetInvoiceRequest{
		Id:        id,
		UserId:    userID,
	})
	if err != nil {
		return nil, err
	}

	invoice = billingmodel.InvoiceFromProto(resp.Invoice)

	return
}

func (s *Controller) GetPayment(ctx context.Context, id, userID int64) (payment *billingmodel.Payment, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("get payment", "id", id, "user_id", userID)

	resp, err := s.client.GetPayment(ctx, &api.GetPaymentRequest{
		Id:       id,
		UserId:   userID,
	})
	if err != nil {
		return nil, err
	}

	payment = billingmodel.PaymentFromProto(resp.Payment)

	return
}
