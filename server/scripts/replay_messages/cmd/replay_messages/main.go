package main

import (
	"os"
	"fmt"
	"time"
	"context"
	"log/slog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/jetstream"
	threadspkg "github.com/bd878/gallery/server/threads/pkg"
)

func main() {
	nc, err := nats.Connect("192.168.10.11:4222")
	if err != nil {
		panic(err)
	}

	rawJs, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	_, err = rawJs.AddStream(&nats.StreamConfig{
		Name: "gallery",
		Subjects: []string{fmt.Sprintf("%s.>", "gallery")},
	})
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	js := jetstream.NewStream("gallery", rawJs)

	pool, err := pgxpool.New(context.TODO(), os.Getenv("PG_CONN"))
	if err != nil {
		panic(err)
	}

	commandStream := am.NewCommandStream(js)

	rows, err := pool.Query(context.TODO(), `
SELECT o.id,o.published_at,o.name,o.subject,s.data
FROM messages_stream.outbox o
INNER JOIN messages_stream.sagas s
	ON o.id = s.id
	AND s.done = false
ORDER BY o.published_at ASC;
	`)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var (
			id, name, subject string
			publishedAt time.Time
			data []byte
		)

		err = rows.Scan(&id, &publishedAt, &name, &subject, &data)
		if err != nil {
			panic(err)
		}

		var msg threads.CreateThread
		err = proto.Unmarshal(data, &msg)
		if err != nil {
			panic(err)
		}

		slog.Debug("msg", slog.Int64("thread_id", msg.ThreadId), slog.String("id", id), slog.String("published_at", publishedAt.String()),
			slog.String("name", name), slog.String("subject", subject))

		cmd := am.NewCommand(threadspkg.CreateThreadCommand, threadspkg.CommandChannel, data)

		slog.Debug("cmd", slog.String("destination", cmd.Destination()), slog.String("id", cmd.ID()))

		cmd.Metadata().Set("COMMAND_REPLY_CHANNEL", "gallery.messages.replies.CreateMessage.prod-3")
		cmd.Metadata().Set("COMMAND_SAGA_ID", id)
		cmd.Metadata().Set("COMMAND_SAGA_NAME", "messages.CreateMessage")

		err = commandStream.Publish(context.TODO(), cmd.Destination(), cmd)
		if err != nil {
			panic(err)
		}
	}
}