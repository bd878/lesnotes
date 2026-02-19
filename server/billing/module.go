package billing

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"context"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/waiter"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/discovery/serf"
	"github.com/bd878/gallery/server/distributed/raft"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/billing/internal/repository/postgres"
	"github.com/bd878/gallery/server/billing/internal/machine"
	"github.com/bd878/gallery/server/billing/internal/controller/application"
	"github.com/bd878/gallery/server/billing/internal/handler/grpc"
)

func Root(ctx context.Context, svc system.Service) (err error) {
	paymentsRepo := postgres.NewPaymentsRepository(svc.Pool(), "billing.payments")
	invoicesRepo := postgres.NewInvoicesRepository(svc.Pool(), "billing.invoices")

	handler := grpc.New(controller)

	if err := setupRaft(svc, paymentsRepo, invoicesRepo); err != nil {
		return err
	}

	if err := setupSerf(svc.Waiter().Context(), svc.Config(), handler); err != nil {
		return err
	}

	api.RegisterBillingServer(svc.RPC(), handler)

	return nil
}

func setupSerf(ctx context.Context, cfg config.Config, handler serf.Handler, logger *logger.Logger) error {
	membership, err := serf.New(
		serf.Config{
			NodeName: cfg.NodeName,
			BindAddr: cfg.SerfAddr,
			Tags: map[string]string{
				"raft_addr": cfg.Addr,
			},
			SerfJoinAddrs: cfg.SerfJoinAddrs,
		},
		handler,
	)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := membership.Leave(); err != nil {
				fmt.Fprintf(os.Stderr, err)
			}
		}()

		err := membership.Run(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, err)
		}
	}()

	return nil
}

func setupRaft(svc system.Service, paymentsRepo *postgres.PaymentsRepository, invoicesRepo *postgres.InvoicesRepository) error {
	raftListener := svc.Mux().Match(func(r io.Reader) bool {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(raft.RaftRPC)}) == 0
	})

	fsm := machine.New(paymentsRepo, invoicesRepo, svc.Logger())

	distributed, err := raft.New(raft.Config{
		Bootstrap:      svc.Config().Bootstrap,
		NodeName:       svc.Config().NodeName,
		RaftLogLevel:   svc.Config().RaftLogLevel,
		DataDir:        svc.Config().DataPath,
		Servers:        svc.Config().RaftServers,
	}, raft.NewStreamLayer(raftListener), fsm, svc.Logger())
	if err != nil {
		return err
	}

	controller := application.New(distributed, paymentsRepo, invoicesRepo, svc.Logger())

	svc.Waiter().Add(func(ctx context.Context) error {
		return distributed.Leave(distributed.NodeName())
	})

	return nil
}
