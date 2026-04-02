package main

import (
	"os"
	"fmt"
	"flag"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	natsLib "github.com/nats-io/nats.go"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/messages/pkg/events"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s nats_addr pg_conn\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	nc, err := natsLib.Connect(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	stream := nats.NewStream(nc)

	fmt.Fprintln(os.Stdout, "nats ok")

	pool, err := pgxpool.New(context.Background(), flag.Arg(1))
	defer pool.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "pg pool ok")

	fmt.Fprintln(os.Stdout, "=== emitting messages private/publish events ===")
	fmt.Fprintln(os.Stdout, "PGConn:", flag.Arg(1), "nats addr:", flag.Arg(0))

	query := "SELECT id, user_id, private FROM messages.messages"

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	messages := make([]*model.Message, 0)
	for rows.Next() {
		message := &model.Message{}

		err := rows.Scan(&message.ID, &message.UserID, &message.Private)
		if err != nil {
			panic(err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	for _, message := range messages {
		if message.Private {
			data, err := proto.Marshal(&api.MessagesPrivated{
				Ids: []int64{message.ID},
				UserId: message.UserID,
				UpdatedAt: message.UpdatedAt,
			})
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(os.Stdout, "private", message.ID)

			err = stream.Publish(context.Background(), events.MessagesChannel, am.NewRawMessage(fmt.Sprintf("%d", message.ID), events.MessagesPrivateEvent, data))
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			}
		} else {
			data, err := proto.Marshal(&api.MessagesPublished{
				Ids:  []int64{message.ID},
				UserId: message.UserID,
				UpdatedAt: message.UpdatedAt,
			})
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(os.Stdout, "publish", message.ID)

			err = stream.Publish(context.Background(), events.MessagesChannel, am.NewRawMessage(fmt.Sprintf("%d", message.ID), events.MessagesPublishEvent, data))
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			}
		}
	}
}