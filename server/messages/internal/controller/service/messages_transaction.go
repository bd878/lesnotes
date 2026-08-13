package service

import (
	"context"
	"log/slog"
	"github.com/jackc/pgx/v5"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type messagesControllerTx struct {
	container di.Container
	MessagesController
}

func NewMessagesControllerTx(container di.Container, controller MessagesController) messagesControllerTx {
	return messagesControllerTx{
		container: container,
		MessagesController: controller,
	}
}

func (s messagesControllerTx) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID, userID int64, private bool, name string) (message *model.Message, err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.SaveMessage(ctx, id, text, title, fileIDs, threadID, userID, private, name)
}

func (s messagesControllerTx) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.DeleteUserMessages(ctx, userID) 
}

func (s messagesControllerTx) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.DeleteMessages(ctx, ids, userID)
}

func (s messagesControllerTx) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.PublishMessages(ctx, ids, userID)
}

func (s messagesControllerTx) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.PrivateMessages(ctx, ids, userID)
}

func (s messagesControllerTx) UpdateMessage(ctx context.Context, id int64, text, title, name *string, fileIDs []int64, userID int64) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.MessagesController.UpdateMessage(ctx, id, text, title, name, fileIDs, userID)
}

func (s messagesControllerTx) closeTx(tx pgx.Tx, err error) error {
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