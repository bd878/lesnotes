package stream

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/am"
	messageevents "github.com/bd878/gallery/server/messages/pkg/events"
	threadsevents "github.com/bd878/gallery/server/threads/pkg/events"
	filesevents "github.com/bd878/gallery/server/files/pkg/events"
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

	err = js.Subscribe(messageevents.MessagesChannel, evtMsgHandler, am.GroupName("search-messages"))
	if err != nil {
		return
	}

	err = js.Subscribe(threadsevents.ThreadsChannel, evtMsgHandler, am.GroupName("search-threads"))
	if err != nil {
		return
	}

	err = js.Subscribe(filesevents.FilesChannel, evtMsgHandler, am.GroupName("search-files"))
	if err != nil {
		return
	}

	err = js.Subscribe(messageevents.TranslationsChannel, evtMsgHandler, am.GroupName("search-translations"))

	return
}