package postgres

import (
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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

func (r *SessionsRepository) Save(ctx context.Context, userID int64, token, createdAt, expiresAt string) (err error) {
	const query = "INSERT INTO %s(user_id, value, created_at, expires_at) VALUES ($1, $2, $3, $4)"

	_, err = r.pool.Exec(ctx, r.table(query), userID, token, createdAt, expiresAt)

	return err
}

func (r *SessionsRepository) Get(ctx context.Context, token string) (session *model.Session, err error) {
	const query = "SELECT user_id, created_at, expires_at FROM %s WHERE value = $1"

	var createdAt, expiresAt *time.Time

	session = &model.Session{
		Token:          token,
	}

	err = r.pool.QueryRow(ctx, r.table(query), token).Scan(&session.UserID, &createdAt, &expiresAt)
	if err != nil {
		return
	}

	session.CreatedAt = createdAt.Format(time.RFC3339)
	session.ExpiresAt = expiresAt.Format(time.RFC3339)

	return
}

func (r *SessionsRepository) List(ctx context.Context, userID int64) (sessions []*model.Session, err error) {
	const query = "SELECT value, created_at, expires_at FROM %s WHERE user_id = $1"

	var rows pgx.Rows
	rows, err = r.pool.Query(ctx, r.table(query), userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions = make([]*model.Session, 0)
	for rows.Next() {
		var createdAt, expiresAt *time.Time

		session := &model.Session{
			UserID:         userID,
		}

		err = rows.Scan(&session.Token, &createdAt, &expiresAt)
		if err != nil {
			return
		}

		session.CreatedAt = createdAt.Format(time.RFC3339)
		session.ExpiresAt = expiresAt.Format(time.RFC3339)

		sessions = append(sessions, session)
	}

	err = rows.Err()

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

func (r SessionsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
