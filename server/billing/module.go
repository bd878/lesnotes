package billing

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/discovery/serf"
	"github.com/bd878/gallery/server/internal/consensus/raft"
	"github.com/bd878/gallery/server/billing/config"
	"github.com/bd878/gallery/server/billing/internal/repository/postgres"
	"github.com/bd878/gallery/server/billing/internal/machine"
	"github.com/bd878/gallery/server/billing/internal/controller/distributed"
	"github.com/bd878/gallery/server/billing/internal/handler/grpc"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	paymentsRepo := postgres.NewPaymentsRepository(svc.Pool(), "billing.payments")
	invoicesRepo := postgres.NewInvoicesRepository(svc.Pool(), "billing.invoices")

	consensus, err := setupRaft(svc, cfg, paymentsRepo, invoicesRepo)
	if err != nil {
		return err
	}

	if err := setupSerf(svc.Waiter().Context(), cfg, consensus, svc.Logger()); err != nil {
		return err
	}

	controller := application.New(consensus, paymentsRepo, invoicesRepo, svc.Logger())

	handler := grpc.New(controller)

	api.RegisterBillingServer(svc.RPC(), handler)

	return nil
}

func setupSerf(ctx context.Context, cfg config.Config, handler serf.Handler, logger *logger.Logger) error {
	membership, err := serf.New(
		serf.Config{
			NodeName: cfg.NodeName,
			BindAddr: cfg.SerfAddr,
			Tags: map[string]string{
				"raft_addr": cfg.RpcAddr,
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
				fmt.Fprintln(os.Stderr, err)
			}
		}()

		membership.Run(ctx)
	}()

	return nil
}

func setupRaft(svc system.Service, cfg config.Config, paymentsRepo *postgres.PaymentsRepository, invoicesRepo *postgres.InvoicesRepository) (*raft.Distributed, error) {
	raftListener := svc.Mux().Match(func(r io.Reader) bool {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(raft.RaftRPC)}) == 0
	})

	fsm := machine.New(paymentsRepo, invoicesRepo, svc.Logger())

	consensus, err := raft.New(raft.Config{
		Bootstrap:      cfg.RaftBootstrap,
		NodeName:       cfg.NodeName,
		RaftLogLevel:   cfg.RaftLogLevel,
		DataDir:        cfg.DataPath,
		Servers:        cfg.RaftServers,
	}, raft.NewStreamLayer(raftListener), fsm, svc.Logger())
	if err != nil {
		return nil, err
	}

	svc.Waiter().Add(func(ctx context.Context) error {
		return consensus.Leave(consensus.NodeName())
	})

	return consensus, nil
}
