package service

import (
	"fmt"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
)

type translationsControllerTx struct {
	container di.Container
	TranslationsController
}

func NewTranslationsControllerTx(container di.Container, controller TranslationsController) translationsControllerTx {
	return translationsControllerTx{
		container:               container,
		TranslationsController: controller,
	}
}

func (s translationsControllerTx) SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit save translation transaction: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.TranslationsController.SaveTranslation(ctx, userID, messageID, lang, title, text)
}

func (s translationsControllerTx) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit update translation transaction: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.TranslationsController.UpdateTranslation(ctx, messageID, lang, title, text)
}

func (s translationsControllerTx) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	ctx = s.container.Scoped(ctx)
	defer func(tx pgx.Tx) {
		err = s.closeTx(tx, err)
		slog.Debug(fmt.Sprintf("commit delete translation transaction: %v", err))
	}(di.Get(ctx, "tx").(pgx.Tx))

	return s.TranslationsController.DeleteTranslation(ctx, messageID, lang)
}

func (s translationsControllerTx) closeTx(tx pgx.Tx, err error) error {
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
