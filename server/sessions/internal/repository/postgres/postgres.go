package repository

import (
	"io"
	"fmt"
	"time"
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

func New(pool *pgxpool.Pool, tableName string) *Repository {
	return &Repository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *Repository) Save(ctx context.Context, userID int64, token string, expiresUTCNano int64) (err error) {
	const query = "INSERT INTO %s(user_id, value, expires_at) VALUES ($1, $2, $3)"

	expiresAt := time.Unix(0, expiresUTCNano)

	_, err = r.pool.Exec(ctx, r.table(query), userID, token, expiresAt)

	return err
}

func (r *Repository) Get(ctx context.Context, token string) (session *model.Session, err error) {
	const query = "SELECT user_id, expires_at FROM %s WHERE value = $1"

	var (
		expiresAt time.Time
		userID int64
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

func (r *Repository) List(ctx context.Context, userID int64) (sessions []*model.Session, err error) {
	const query = "SELECT value, expires_at FROM %s WHERE user_id = $1"

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
	const query = "DELETE FROM %s WHERE value = $1"

	_, err = r.pool.Exec(ctx, r.table(query), token)

	return
}

func (r *Repository) DeleteAll(ctx context.Context, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE user_id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), userID)

	return
}

func (r *Repository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.table("TRUNCATE TABLE %s"))
	return
}

func (r *Repository) Dump(ctx context.Context) (reader io.ReadCloser, err error) {
	var (
		writer io.WriteCloser
		conn   *pgxpool.Conn
	)

	query := r.table("COPY %s TO STDOUT BINARY")

	reader, writer = io.Pipe()

	conn, err = r.pool.Acquire(ctx)
	if err != nil {
		conn.Release()
		return
	}

	go func(ctx context.Context, query string, conn *pgxpool.Conn, writer io.WriteCloser) {
		_, err := conn.Conn().PgConn().CopyTo(ctx, writer, query)
		defer writer.Close()
		defer conn.Release()
		if err != nil {
			logger.Errorw("failed to dump", "error", err)
		}
	}(ctx, query, conn, writer)

	return
}

func (r *Repository) Restore(ctx context.Context, reader io.ReadCloser) (err error) {
	var conn *pgxpool.Conn

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	if err != nil {
		conn.Release()
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)
	defer conn.Release()

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
