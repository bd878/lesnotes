package service

import (
	"fmt"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
)

type usersControllerTx struct {
	container di.Container
	Controller
}

func NewUsersControllerTx(container di.Container, controller Controller) usersControllerTx {
	return usersControllerTx{
		container:  container,
		Controller: controller,
	}
}

func (s usersControllerTx) DeleteUser(ctx context.Context, id int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit delete user transaction: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.DeleteUser(ctx, id)
}

func (s usersControllerTx) closeTx(tx pgx.Tx, err error) error {
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
