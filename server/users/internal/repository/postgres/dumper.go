package postgres

import (
	"fmt"
	"time"
	"sync"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type UsersDumper struct {
	pool              *pgxpool.Pool
	usersTableName    string
	premiumsTableName string
	rows              pgx.Rows
	ctx               context.Context
	cancel            context.CancelCauseFunc
	ch                chan *api.UserSnapshot
	wg                sync.WaitGroup
}

func NewUsersDumper(pool *pgxpool.Pool, usersTableName, premiumsTableName string) *UsersDumper {
	return &UsersDumper{
		pool:                 pool,
		usersTableName:       usersTableName,
		premiumsTableName:    premiumsTableName,
	}
}

func (r *UsersDumper) Open(ctx context.Context) (ch chan *api.UserSnapshot, err error) {
	query := "SELECT id, login, salt, metadata, created_at, updated_at FROM %s"

	r.ctx, r.cancel = context.WithCancelCause(ctx)
	r.rows, err = r.pool.Query(r.ctx, r.usersTable(query))
	if err != nil {
		return
	}

	ch = make(chan *api.UserSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go func() {
		defer close(r.ch)
		defer r.rows.Close()
		defer r.wg.Done()

		for r.rows.Next() {
			user := &api.UserSnapshot{}

			var createdAt, updatedAt *time.Time
			err = r.rows.Scan(&user.Id, &user.Login, &user.HashedPassword, &user.Metadata, &createdAt, &updatedAt)
			if err != nil {
				logger.Errorln(err)
				r.cancel(err)
				return
			}

			user.CreatedAt = createdAt.Format(time.RFC3339)
			user.UpdatedAt = updatedAt.Format(time.RFC3339)

			select {
			case <-r.ctx.Done():
				return
			default:
			}

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
	logger.Debugln("close dumper")
	r.cancel(nil)
	r.wg.Wait()
	return nil
}

func (r *UsersDumper) Restore(ctx context.Context, user *api.UserSnapshot) (err error) {
	query := "INSERT INTO %s(id, login, salt, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.usersTable(query), user.Id, user.Login, user.HashedPassword, user.Metadata, user.CreatedAt, user.UpdatedAt)

	return
}

func (r UsersDumper) usersTable(query string) string {
	return fmt.Sprintf(query, r.usersTableName)
}
