package service

import (
	"fmt"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
)

type billingControllerTx struct {
	container di.Container
	Controller
}

func NewBillingControllerTx(container di.Container, controller Controller) billingControllerTx {
	return billingControllerTx{
		container:  container,
		Controller: controller,
	}
}

func (s billingControllerTx) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit proceed payment transaction: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.ProceedPayment(ctx, id, userID)
}

func (s billingControllerTx) closeTx(tx pgx.Tx, err error) error {
	ctx := context.Background()

	p := recover()
	switch {
	case p != nil:
		_ = tx.Rollback(ctx)
		panic(p)
	case err != nil:
		slog.Error("rollback with error", slog.String("error", err.Error()))
		err = tx.Rollback(ctx)
		return err
	default:
		return tx.Commit(ctx)
	}
}
