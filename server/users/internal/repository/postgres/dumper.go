package postgres

import (
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type UsersDumper struct {
	pool       *pgxpool.Pool
	tableName  string
	rows       pgx.Rows
	ctx        context.Context
	cancel     context.CancelCauseFunc
	ch         chan *model.User
}

func NewUsersDumper(pool *pgxpool.Pool, tableName string) *UsersDumper {
	return &UsersDumper{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *UsersDumper) Open(ctx context.Context) (ch chan *model.User, err error) {
	query := "SELECT id, login, salt, metadata, created_at, updated_at FROM %s"

	r.ctx, r.cancel = context.WithCancelCause(ctx)
	r.rows, err = r.pool.Query(r.ctx, r.table(query))
	if err != nil {
		return
	}

	ch = make(chan *model.User, 100)
	r.ch = ch

	go func() {
		defer close(r.ch)
		for r.rows.Next() {
			user := &model.User{}

			var createdAt, updatedAt *time.Time
			err = r.rows.Scan(&user.ID, &user.Login, &user.HashedPassword, &user.Metadata, &createdAt, &updatedAt)
			if err != nil {
				logger.Errorln(err)
				r.cancel(err)
				return
			}

			user.CreatedAt = createdAt.Format(time.RFC3339)
			user.UpdatedAt = updatedAt.Format(time.RFC3339)

			r.ch <- user
		}

		if err := r.rows.Err(); err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}
	}()

	return
}

func (r *UsersDumper) Close() (err error) {
	r.rows.Close()
	return nil
}

func (r *UsersDumper) Restore(ctx context.Context, user *model.User) (err error) {
	query := "INSERT INTO %s(id, login, salt, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.table(query), user.ID, user.Login, user.HashedPassword, user.Metadata, user.CreatedAt, user.UpdatedAt)

	return
}

func (r UsersDumper) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
