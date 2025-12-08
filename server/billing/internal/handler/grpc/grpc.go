package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/billing/pkg/model"
)

type Controller interface {
	GetServers(ctx context.Context) (servers []*api.Server, err error)
	CreateInvoice(ctx context.Context, id string, userID int64, currency string, total int64, metadata []byte) (err error)
	StartPayment(ctx context.Context, id, userID int64, invoiceID string, currency string, total int64, metadata []byte) (err error)
	ProceedPayment(ctx context.Context, id, userID int64) (err error)
	CancelPayment(ctx context.Context, id, userID int64) (err error)
	RefundPayment(ctx context.Context, id, userID int64) (err error)
	GetInvoice(ctx context.Context, id, userID int64) (invoice *model.Invoice, err error)
	GetPayment(ctx context.Context, id, userID int64) (payment *model.Payment, err error)
}

type Handler struct {
	api.UnimplementedBillingServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
}

func (h *Handler) CreateInvoice(ctx context.Context, req *api.CreateInvoiceRequest) (resp *api.CreateInvoiceResponse, err error) {
	err := h.controller.CreateInvoice(ctx, req.Id, req.UserId, req.Currency, req.Total, req.Metadata)
	if err != nil {
		return nil, err
	}

	resp = &api.CreateInvoiceResponse{
	}

	return
}

func (h *Handler) StartPayment(ctx context.Context, req *api.StartPaymentRequest) (resp *api.StartPaymentResponse, err error) {
	err := h.controller.StartPayment(ctx, req.Id, req.UserId, req.InvoiceId, req.Currency, req.Total, req.Metadata)
	if err != nil {
		return nil, err
	}

	resp = &api.StartPaymentResponse{
	}

	return
}

func (h *Handler) ProceedPayment(ctx context.Context, req *api.ProceedPaymentRequest) (resp *api.ProceedPaymentResponse, err error) {
	err := h.controller.ProceedPayment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.ProceedPaymentResponse{
	}

	return
}

func (h *Handler) CancelPayment(ctx context.Context, req *api.CancelPaymentRequest) (resp *api.CancelPaymentResponse, err error) {
	err := h.controller.CancelPayment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.CancelPaymentResponse{
	}

	return
}

func (h *Handler) RefundPayment(ctx context.Context, req *api.RefundPaymentRequest) (resp *api.RefundPaymentResponse, err error) {
	err := h.controller.RefundPayment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.RefundPaymentResponse{
	}

	return
}

func (h *Handler) GetInvoice(ctx context.Context, req *api.GetInvoiceRequest) (resp *api.GetInvoiceResponse, err error) {
	invoice, err := h.controller.GetInvoice(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.GetInvoiceResponse{
		Invoice: model.InvoiceToProto(invoice),
	}

	return
}

func (h *Handler) GetPayment(ctx context.Context, req *api.GetPaymentRequest) (resp *api.GetPaymentResponse, err error) {
	payment, err := h.controller.GetPayment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.GetPaymentResponse{
		Invoice: model.PaymentToProto(payment),
	}

	return
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (resp *api.GetServersResponse, err error) {
	servers, err := h.controller.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	resp = &api.GetServersResponse{
		Servers: servers,
	}

	return
}
