package postgres

import (
	"io"
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

type SessionsRepository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func NewSessionsRepository(pool *pgxpool.Pool, tableName string) *SessionsRepository {
	return &SessionsRepository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *SessionsRepository) Save(ctx context.Context, userID int64, token string, expiresUTCNano int64) (err error) {
	const query = "INSERT INTO %s(user_id, value, expires_at) VALUES ($1, $2, $3)"

	expiresAt := time.Unix(0, expiresUTCNano)

	_, err = r.pool.Exec(ctx, r.table(query), userID, token, expiresAt)

	return err
}

func (r *SessionsRepository) Get(ctx context.Context, token string) (session *model.Session, err error) {
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

func (r *SessionsRepository) List(ctx context.Context, userID int64) (sessions []*model.Session, err error) {
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

func (r *SessionsRepository) Delete(ctx context.Context, token string) (err error) {
	const query = "DELETE FROM %s WHERE value = $1"

	_, err = r.pool.Exec(ctx, r.table(query), token)

	return
}

func (r *SessionsRepository) DeleteAll(ctx context.Context, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE user_id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), userID)

	return
}

func (r *SessionsRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping sessions repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *SessionsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring sessions repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r SessionsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
