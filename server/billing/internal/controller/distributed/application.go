package application

import (
	"time"
	"context"
	"bytes"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/billing/pkg/model"
	"github.com/bd878/gallery/server/billing/internal/machine"
)

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

func (m *Distributed) Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), duration)
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