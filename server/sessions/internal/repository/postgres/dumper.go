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

type Dumper struct {
	pool              *pgxpool.Pool
	sessionsTableName string
	ctx               context.Context
	cancel            context.CancelCauseFunc
	ch                chan *api.SessionsSnapshot
	wg                sync.WaitGroup
}

func NewDumper(pool *pgxpool.Pool, sessionsTableName string) *Dumper {
	return &Dumper{
		pool:                 pool,
		sessionsTableName:    sessionsTableName,
	}
}

func (r *Dumper) Open(ctx context.Context) (ch chan *api.SessionsSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.SessionsSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.runSessions()

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()

	return
}

func (r *Dumper) runSessions() {
	query := "SELECT user_id, value, created_at, expires_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("sessions dump finished")

	rows, err := r.pool.Query(r.ctx, r.sessionsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		session := &api.SessionSnapshotItem{}

		var createdAt, expiresAt *time.Time
		err = rows.Scan(&session.UserId, &session.Token, &createdAt, &expiresAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		session.CreatedAt = createdAt.Format(time.RFC3339)
		session.ExpiresAt = expiresAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.SessionsSnapshot{
			Item: &api.SessionsSnapshot_Session{
				Session: session,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) Close() (err error) {
	logger.Debugln("close dumper")
	r.cancel(nil)
	r.wg.Wait()
	return nil
}

func (r *Dumper) Restore(ctx context.Context, snapshot *api.SessionsSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.SessionsSnapshot_Session:

		query := "INSERT INTO %s(user_id, value, created_at, expires_at) VALUES ($1, $2, $3, $4)"

		_, err = r.pool.Exec(ctx, r.sessionsTable(query), v.Session.UserId, v.Session.Token, v.Session.CreatedAt, v.Session.ExpiresAt)

		return

	default:
		return errors.New("unknown snapshot item")
	}
}

func (r Dumper) sessionsTable(query string) string {
	return fmt.Sprintf(query, r.sessionsTableName)
}

