package saga

import (
	"context"
	"log/slog"
	"github.com/jackc/pgx/v5"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
)

func RegisterReplyHandlersTx(c di.Container) error {
	replyMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.RawMessage) (err error) {
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

		slog.Debug("handle reply",
			slog.String("name", msg.MessageName()),
			slog.String("subject", msg.Subject()),
		)

		replyHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewReplyMessageHandler(
				di.Get(ctx, "replyEventHandlers").(ddd.ReplyHandler[am.ReplyMessage]),
			),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return replyHandlers.HandleMessage(ctx, msg)
	})

	js := c.Get("js").(am.RawMessageStream)

	return js.Subscribe(messageReplyChannel + ".>", replyMsgHandler)
}