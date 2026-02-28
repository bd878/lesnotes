package main

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/bd878/gallery/server/threads/pkg/loadbalance"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/threads/pkg/model"
	"github.com/bd878/gallery/server/threads/pkg/loadbalance"
)

func main() {
	PGConn, ok := os.LookupEnv("PG_CONN")
	if !ok {
		panic("PG_CONN env required")
	}
	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		panic("ADDR env required")
	}

	fmt.Fprintln(os.Stdout, "=== running threads migration ===")
	fmt.Fprintln(os.Stdout, "PGConn:", PGConn, "addr:", addr)

	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			addr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "grpc conn ok")

	client := api.NewThreadsClient(conn)
	pool, err := pgxpool.New(context.Background(), PGConn)
	defer pool.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "pg pool ok")

	query := "SELECT id, name, private, user_id, parent_id, next_id, prev_id, description FROM threads.threads"

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "rows ok")

	threads := make([]*model.Thread, 0)
	for rows.Next() {
		thread := &model.Thread{}

		err := rows.Scan(&thread.ID, &thread.Name, &thread.Private, &thread.UserID, &thread.ParentID, &thread.NextID, &thread.PrevID, &thread.Description)
		if err != nil {
			panic(err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "scanned threads rows:", len(threads))

	for _, thread := range threads {
		fmt.Fprintln(os.Stdout, "create thread", "id", thread.ID)
		_, err = client.Create(context.Background(), &api.CreateRequest{
			Id:                  thread.ID,
			UserId:              thread.UserID,
			ParentId:            thread.ParentID,
			Name:                thread.Name,
			Private:             thread.Private,
			NextId:              thread.NextID,
			PrevId:              thread.PrevID,
			Description:         thread.Description,
		})
		if err != nil {
			fmt.Fprintln(os.Stdout, "[WARN]: failed to save thread:", err)
		}
	}

	fmt.Fprintln(os.Stdout, "=== threads migration done ===")
}