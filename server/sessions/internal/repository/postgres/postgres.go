package repository

import (
	"fmt"
	"os"
	"io"
	"time"
	"errors"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		tableName: "sessions.sessions",
		pool:      pool,
	}
}

func (r *Repository) Save(ctx context.Context, userID int32, token string, expiresUTCNano int64) (err error) {
	const query = "INSERT INTO %s(user_id, value, expires_at) VALUES ($1, $2, $3)"

	expiresAt := time.Unix(0, expiresUTCNano)

	_, err = r.pool.Exec(ctx, r.table(query), userID, token, expiresAt)

	return err
}

func (r *Repository) Get(ctx context.Context, token string) (session *model.Session, err error) {
	const query := "SELECT user_id, expires_at FROM %s WHERE value = $1"

	var (
		expiresAt time.Time
		userID int32
	)

	err = r.pool.QueryRow(ctx, r.table(query), token).Scan(&userID, &expiresAt)
	if err != nil {
		return
	}

	session = &model.Session{
		UserID:         userID,
		Token:          token,
		ExpiresUTCNano: expiresAt.UnixNano(),
	}

	return
}

func (r *Repository) List(ctx context.Context, userID int32) (sessions []*model.Session, err error) {
	const query := "SELECT value, expires_at FROM %s WHERE user_id = $1"

	var rows pgx.Rows
	rows, err = r.pool.Query(ctx, r.table(query), userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions = make([]*model.Session, 0)
	for rows.Next() {
		var (
			token string
			expiresAt time.Time
		)

		err = rows.Scan(&token, &expiresAt)
		if err != nil {
			return
		}

		sessions = append(sessions, &model.Session{
			UserID:         userID,
			Token:          token,
			ExpiresUTCNano: expiresAt.UnixNano(),
		})
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *Repository) Delete(ctx context.Context, token string) (err error) {
	const query := "DELETE FROM %s WHERE value = $1"

	_, err = r.pool.Exec(ctx, r.table(query), token)

	return
}

func (r *Repository) DeleteAll(ctx context.Context, userID int32) (err error) {
	const query := "DELETE FROM %s WHERE user_id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), userID)

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
