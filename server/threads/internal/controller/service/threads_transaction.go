package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
)

type threadsControllerTx struct {
	container di.Container
	*Controller
}

func NewThreadsControllerTx(container di.Container, controller *Controller) *threadsControllerTx {
	return &threadsControllerTx{
		container:  container,
		Controller: controller,
	}
}

func (s *threadsControllerTx) PublishThread(ctx context.Context, id, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit publish thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.PublishThread(ctx, id, userID)
}

func (s *threadsControllerTx) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit private thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.PrivateThread(ctx, id, userID)
}

func (s *threadsControllerTx) CreateThread(ctx context.Context, id, userID, parentID int64, name, description, title string, private bool) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit create thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.CreateThread(ctx, id, userID, parentID, name, description, title, private)
}

func (s *threadsControllerTx) UpdateThread(ctx context.Context, id, userID int64, name, description, title *string) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit update thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.UpdateThread(ctx, id, userID, name, description, title)
}

func (s *threadsControllerTx) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit reorder thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.ReorderThread(ctx, id, userID, parentID, nextID, prevID)
}

func (s *threadsControllerTx) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit private messages transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.PrivateMessages(ctx, ids, userID)
}

func (s *threadsControllerTx) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit publish messages transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.PublishMessages(ctx, ids, userID)
}

func (s *threadsControllerTx) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit delete thread transaction, err: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.Controller.DeleteThread(ctx, id, userID)
}

func (s *threadsControllerTx) closeTx(tx pgx.Tx, err error) error {
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