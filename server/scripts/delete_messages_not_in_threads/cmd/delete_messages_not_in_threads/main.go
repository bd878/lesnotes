package main

import (
	"os"
	"fmt"
	"context"
	"log/slog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/internal/jetstream"
	messagespkg "github.com/bd878/gallery/server/messages/pkg"
)

func main() {
	nc, err := nats.Connect("192.168.10.111:4222")
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
SELECT id, user_id FROM messages.messages m
WHERE NOT EXISTS (
	SELECT 1 FROM threads.threads t WHERE t.id = m.id
);
	`)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	list := make([]*messages.DeleteCommand, 0)
	for rows.Next() {
		cmd := &messages.DeleteCommand{}
		err := rows.Scan(&cmd.Id, &cmd.UserId)
		if err != nil {
			panic(err)
		}

		list = append(list, cmd)
	}

	for _, cmd := range list {
		data, err := proto.Marshal(cmd)
		if err != nil {
			panic(err)
		}

		slog.Debug("msg", slog.Int64("id", cmd.Id), slog.Int64("user_id", cmd.UserId))

		msg := am.NewCommand(messagespkg.DeleteMessageCommand, messagespkg.CommandChannel, data)

		slog.Debug("cmd", slog.String("destination", msg.Destination()), slog.String("id", msg.ID()))

		msg.Metadata().Set("COMMAND_REPLY_CHANNEL", "gallery.messages.replies.DeleteMessage.stage-3")
		msg.Metadata().Set("COMMAND_SAGA_ID", msg.ID())
		msg.Metadata().Set("COMMAND_SAGA_NAME", "messages.DeleteMessage")

		err = commandStream.Publish(context.TODO(), msg.Destination(), msg)
		if err != nil {
			panic(err)
		}
	}
}