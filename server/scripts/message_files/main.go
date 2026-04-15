package main

import (
	"os"
	"fmt"
	"flag"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	natsLib "github.com/nats-io/nats.go"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/nats"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/pkg/events"
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

	fmt.Fprintln(os.Stdout, "=== emitting message created with files event ===")
	fmt.Fprintln(os.Stdout, "PGConn:", flag.Arg(1), "nats addr:", flag.Arg(0))

	query := "SELECT user_id, message_id, array_agg(file_id ORDER BY file_id DESC) FROM messages.files GROUP BY user_id, message_id"

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	list := make([]*api.MessageCreated, 0)
	for rows.Next() {
		event := &api.MessageCreated{}

		err := rows.Scan(&event.UserId, &event.Id, &event.FileIds)
		if err != nil {
			panic(err)
		}

		logger.Infoln(event)

		list = append(list, event)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	for _, event := range list {
		data, err := proto.Marshal(event)
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stdout, "publish", event.Id)

		err = stream.Publish(context.Background(), events.MessagesChannel, am.NewRawMessage(fmt.Sprintf("%d", event.Id), events.MessageCreatedEvent, data))
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
		}
	}
}