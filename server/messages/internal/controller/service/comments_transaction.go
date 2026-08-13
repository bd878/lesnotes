package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
)

type commentsControllerTx struct {
	container di.Container
	CommentsController
}

func NewCommentsControllerTx(container di.Container, controller CommentsController) commentsControllerTx {
	return commentsControllerTx{
		container:          container,
		CommentsController: controller,
	}
}

func (s commentsControllerTx) SendComment(ctx context.Context, id, userID, messageID int64, text string, metadata []byte) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug("commit send comment transaction", slog.String("err", err.Error()))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.CommentsController.SendComment(ctx, id, userID, messageID, text, metadata)
}

func (s commentsControllerTx) UpdateComment(ctx context.Context, id, userID int64, text *string) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug("commit update comment transaction", slog.String("err", err.Error()))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.CommentsController.UpdateComment(ctx, id, userID, text)
}

func (s commentsControllerTx) DeleteComment(ctx context.Context, id, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug("commit delete comment transaction", slog.String("err", err.Error()))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.CommentsController.DeleteComment(ctx, id, userID)
}

func (s commentsControllerTx) DeleteMessageComments(ctx context.Context, messageID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug("commit delete message comments transaction", slog.String("err", err.Error()))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.CommentsController.DeleteMessageComments(ctx, messageID)
}

func (s commentsControllerTx) closeTx(tx pgx.Tx, err error) error {
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
