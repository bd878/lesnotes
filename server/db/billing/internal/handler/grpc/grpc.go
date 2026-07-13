package grpc

import (
	"context"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/billing"
	"github.com/bd878/gallery/server/db/billing/pkg/machine"
)

type Controller interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
	GetInvoice(ctx context.Context, id string, userID int64) (invoice *billing.Invoice, err error)
	GetPayment(ctx context.Context, id, userID int64) (payment *billing.Payment, err error)
}

type Handler struct {
	billing.UnimplementedBillingServer
	api.UnimplementedDistributedServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
}

func (h *Handler) Apply(ctx context.Context, req *api.Command) (resp *api.CommandResponse, err error) {
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	err = h.controller.Apply(ctx, machine.RequestType(req.ReqType), req.Cmd, duration)

	resp = &api.CommandResponse{}

	return
}

func (h *Handler) GetInvoice(ctx context.Context, req *billing.GetInvoiceRequest) (resp *billing.GetInvoiceResponse, err error) {
	invoice, err := h.controller.GetInvoice(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &billing.GetInvoiceResponse{
		Invoice: invoice,
	}

	return
}

func (h *Handler) GetPayment(ctx context.Context, req *billing.GetPaymentRequest) (resp *billing.GetPaymentResponse, err error) {
	payment, err := h.controller.GetPayment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &billing.GetPaymentResponse{
		Payment: payment,
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
