package stream

import (
	"context"
	"log/slog"
	"github.com/jackc/pgx/v5"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
)

func RegisterCommandHandlersTx(c di.Container) (err error) {
	cmdMsgHandlers := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.RawMessage) error {
		ctx = c.Scoped(ctx)
		defer func(tx pgx.Tx) {
			p := recover()
			switch {
			case p != nil:
				_ = tx.Rollback(ctx)
				panic(p)
			case err != nil:
				slog.Error("rollback with error", slog.String("error", err.Error()))
				err = tx.Rollback(ctx)
			default:
				err = tx.Commit(ctx)
			}
		}(di.Get(ctx, "tx").(pgx.Tx))

		cmdMsgHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewCommandMessageHandler(
				di.Get(ctx, "replyStream").(am.ReplyStream),
				di.Get(ctx, "commandHandlers").(ddd.CommandHandler[ddd.Command]),
			).(am.RawMessageHandler),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return cmdMsgHandlers.HandleMessage(ctx, msg)
	})

	js := c.Get("js").(am.RawMessageStream)

	return RegisterCommandHandlers(js, cmdMsgHandlers)
}