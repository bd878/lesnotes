package postgres

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type UsersDumper struct {
	pool              *pgxpool.Pool
	usersTableName    string
	premiumsTableName string
	ctx               context.Context
	cancel            context.CancelCauseFunc
	ch                chan *api.UsersSnapshot
	wg                sync.WaitGroup
}

func NewUsersDumper(pool *pgxpool.Pool, usersTableName, premiumsTableName string) *UsersDumper {
	return &UsersDumper{
		pool:                 pool,
		usersTableName:       usersTableName,
		premiumsTableName:    premiumsTableName,
	}
}

func (r *UsersDumper) Open(ctx context.Context) (ch chan *api.UsersSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.UsersSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.runUsers()
	r.wg.Add(1)
	go r.runPremiums()

	return
}

func (r *UsersDumper) runUsers() {
	query := "SELECT id, login, salt, metadata, created_at, updated_at FROM %s"

	defer r.wg.Done()

	rows, err := r.pool.Query(r.ctx, r.usersTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := &api.UserSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&user.Id, &user.Login, &user.HashedPassword, &user.Metadata, &createdAt, &updatedAt)
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

		r.ch <- &api.UsersSnapshot{
			Item: &api.UsersSnapshot_User{
				User: user,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *UsersDumper) runPremiums() {
	query := "SELECT id, invoice_id, created_at, expires_at FROM %s"

	defer r.wg.Done()

	rows, err := r.pool.Query(r.ctx, r.premiumsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		premium := &api.PremiumSnapshotItem{}

		var createdAt, expiresAt *time.Time
		err = rows.Scan(&premium.Id, &premium.InvoiceId, &createdAt, &expiresAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		premium.CreatedAt = createdAt.Format(time.RFC3339)
		premium.ExpiresAt = expiresAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.UsersSnapshot{
			Item: &api.UsersSnapshot_Premium{
				Premium: premium,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *UsersDumper) Close() (err error) {
	logger.Debugln("close dumper")
	r.cancel(nil)
	r.wg.Wait()
	close(r.ch)
	return nil
}

func (r *UsersDumper) Restore(ctx context.Context, snapshot *api.UsersSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.UsersSnapshot_User:

		query := "INSERT INTO %s(id, login, salt, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

		_, err = r.pool.Exec(ctx, r.usersTable(query), v.User.Id, v.User.Login, v.User.HashedPassword, v.User.Metadata, v.User.CreatedAt, v.User.UpdatedAt)

		return

	case *api.UsersSnapshot_Premium:

		query := "INSERT INTO %s(id, invoice_id, created_at, expires_at) VALUES ($1, $2, $3, $4)"

		_, err = r.pool.Exec(ctx, r.premiumsTable(query), v.Premium.Id, v.Premium.InvoiceId, v.Premium.CreatedAt, v.Premium.ExpiresAt)

		return

	default:
		return errors.New("unknown snapshot item")
	}
}

func (r UsersDumper) usersTable(query string) string {
	return fmt.Sprintf(query, r.usersTableName)
}

func (r UsersDumper) premiumsTable(query string) string {
	return fmt.Sprintf(query, r.premiumsTableName)
}
