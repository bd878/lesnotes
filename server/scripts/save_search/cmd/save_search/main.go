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
	"github.com/bd878/gallery/server/api/threads"
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
SELECT id, user_id, parent_id, name, description, private, created_at, updated_at
FROM threads.threads
	`)

	defer rows.Close()
	if err != nil {
		panic(err)
	}

	list := make([]*threads.Thread, 0)
	for rows.Next() {
		thread := &threads.Thread{}

		var createdAt, updatedAt time.Time

		err := rows.Scan(&thread.Id, &thread.UserId, &thread.ParentId, &thread.Name, &thread.Description, &thread.Private, &createdAt, &updatedAt)
		if err != nil {
			panic(err)
		}

		thread.CreatedAt = createdAt.Format(time.RFC3339)
		thread.UpdatedAt = updatedAt.Format(time.RFC3339)

		list = append(list, thread)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, thread := range list {
		cmd, err := proto.Marshal(&search.AppendThreadCommand{
			Id:        thread.Id,
			Description:      thread.Description,
			ParentId:  thread.ParentId,
			Name:      thread.Name,
			UserId:    thread.UserId,
			Private:   thread.Private,
			CreatedAt: thread.CreatedAt,
			UpdatedAt: thread.UpdatedAt,
		})
		if err != nil {
			panic(err)
		}

		_, err = client.Apply(context.TODO(), &api.Command{
			ReqType:  int32(machine.AppendThreadRequest),
			Cmd:      cmd,
			Duration: "10s",
		})
		if err != nil {
			fmt.Fprintf(os.Stdout, "%v\n", err)
		}
	}
}