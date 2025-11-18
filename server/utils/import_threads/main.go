package main

import (
	"os"
	"fmt"
	"flag"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/threads/pkg/loadbalance"
)


func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s addr pg_conn\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			flag.Arg(0),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := api.NewThreadsClient(conn)

	pool, err := pgxpool.New(context.Background(), flag.Arg(1))
	if err != nil {
		panic(err)
	}

	const query = "SELECT id, user_id, thread_id, name, private FROM messages.messages"

	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		panic(err)
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(context.Background())
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(context.Background())
		default:
			err = tx.Commit(context.Background())
		}

		os.Exit(1)
	}()

	rows, err := tx.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var (
			id, userID, threadID int64
			name string
			private bool
		)

		err = rows.Scan(&id, &userID, &threadID, &name, &private)
		if err != nil {
			panic(err)
		}

		_, err = client.Create(context.Background(), &api.CreateRequest{
			Id:       id,
			UserId:   userID,
			ParentId: threadID,
			Name:     name,
			Private:  private,
		})
		if err != nil {
			logger.Errorw("failed to create", "error", err)
		}
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
}