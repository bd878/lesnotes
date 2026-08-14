package stream

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/am"
	messagesevents "github.com/bd878/gallery/server/messages/pkg/events"
)

func RegisterIntegrationEventHandlersTx(c di.Container) (err error) {
	evtMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.RawMessage) error {
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

		evtHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(
				di.Get(ctx, "integrationEventHandlers").(am.MessageHandler[am.EventMessage]),
			),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return evtHandlers.HandleMessage(ctx, msg)
	})

	js := c.Get("js").(am.RawMessageStream)

	return js.Subscribe(messagesevents.MessagesChannel, evtMsgHandler, am.GroupName("files-messages"))
}
