package main

import (
	"os"
	"fmt"
	"time"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/search"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/db/search/pkg/machine"
	"github.com/bd878/gallery/server/db/search/pkg/loadbalance"
)

func main() {
	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			"192.168.10.140:9114",
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := search.NewSearchClient(conn)

	pool, err := pgxpool.New(context.TODO(), os.Getenv("PG_CONN"))
	if err != nil {
		panic(err)
	}

	rows, err := pool.Query(context.TODO(), `
SELECT id, user_id, name, title, text, private, created_at, updated_at
FROM messages.messages
	`)

	defer rows.Close()
	if err != nil {
		panic(err)
	}

	list := make([]*messages.Message, 0)
	for rows.Next() {
		message := &messages.Message{}

		var createdAt, updatedAt time.Time

		err := rows.Scan(&message.Id, &message.UserId, &message.Name, &message.Title, &message.Title, &message.Private, &createdAt, &updatedAt)
		if err != nil {
			panic(err)
		}

		message.CreatedAt = createdAt.Format(time.RFC3339)
		message.UpdatedAt = updatedAt.Format(time.RFC3339)

		list = append(list, message)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, message := range list {
		cmd, err := proto.Marshal(&search.AppendMessageCommand{
			Id:        message.Id,
			Text:      message.Text,
			Title:     message.Title,
			Name:      message.Name,
			UserId:    message.UserId,
			Private:   message.Private,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		})
		if err != nil {
			panic(err)
		}

		_, err = client.Apply(context.TODO(), &api.Command{
			ReqType:  int32(machine.AppendMessageRequest),
			Cmd:      cmd,
			Duration: "10s",
		})
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v\n", err)
		}
	}
}