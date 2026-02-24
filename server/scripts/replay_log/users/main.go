package main

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/pkg/loadbalance"
)

func main() {
	PGConn := os.Getenv("PG_CONN")
	addr := os.Getenv("ADDR")
	tableName := os.Getenv("TABLE")

	fmt.Fprintln(os.Stdout, "=== running users migration ===")
	fmt.Fprintln(os.Stdout, "PGConn", PGConn, "addr", addr, "table", tableName)

	table := func(query string) string {
		return fmt.Sprintf(query, tableName)
	}

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

	client := api.NewUsersClient(conn)
	pool, err := pgxpool.New(context.Background(), PGConn)
	defer pool.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "pg pool ok")

	query := table("SELECT id, login, metadata FROM %s")

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "rows ok")

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}

		err := rows.Scan(&user.ID, &user.Login, &user.Metadata)
		if err != nil {
			panic(err)
		}

		if user.ID != model.PublicUserID {
			users = append(users, user)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "scanned users rows:", len(users))

	for _, user := range users {
		fmt.Fprintln(os.Stdout, "create user", "id", user.ID, "login", user.Login)
		_, err = client.CreateUser(context.Background(), &api.CreateUserRequest{
			Id:           user.ID,
			Login:        user.Login,
			Password:     "12345",
			Metadata:     user.Metadata,
		})
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintln(os.Stdout, "=== users migration done ===")
}