package orchestrator

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
)

func RegisterReplyHandlers(c di.Container) (err error) {
	replyMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.RawMessage) error {
		ctx = c.Scoped(ctx)
		defer func(tx pgx.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback(ctx)
				panic(p)
			} else if err != nil {
				err = tx.Rollback(ctx)
			} else {
				err = tx.Commit(ctx)
			}
		}(di.Get(ctx, "tx").(pgx.Tx))

		replyHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewReplyMessageHandler(
				di.Get(ctx, "replyHandlers").(ddd.ReplyHandler[ddd.Reply]),
			),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return replyHandlers.HandleMessage(ctx, msg)
	})

	js := c.Get("js").(am.RawMessageStream)

	return js.Subscribe(CreateMessageReplyChannel, replyMsgHandler, am.GroupName("messages-replies"))
}

type replyHandlers struct {}

func NewReplyHandlers() replyHandlers {
	return replyHandlers{}
}

func (h replyHandlers) HandleReply(ctx context.Context, reply am.ReplyMessage) error {
	return nil
}