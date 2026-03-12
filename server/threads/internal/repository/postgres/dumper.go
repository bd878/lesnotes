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
	threadsTableName  string
	ctx               context.Context
	cancel            context.CancelCauseFunc
	ch                chan *api.ThreadsSnapshot
	wg                sync.WaitGroup
}

func NewDumper(pool *pgxpool.Pool, threadsTableName string) *Dumper {
	return &Dumper{
		pool:                 pool,
		threadsTableName:     threadsTableName,
	}
}

func (r *Dumper) Open(ctx context.Context) (ch chan *api.ThreadsSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.ThreadsSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.runThreads()

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()

	return
}

func (r *Dumper) runThreads() {
	query := "SELECT id, name, private, user_id, parent_id, next_id, prev_id, created_at, updated_at, description FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("threads dump finished")

	rows, err := r.pool.Query(r.ctx, r.threadsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		thread := &api.ThreadSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&thread.Id, &thread.Name, &thread.Private, &thread.UserId, &thread.ParentId,
			&thread.NextId, &thread.PrevId, &createdAt, &updatedAt, &thread.Description)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		thread.CreatedAt = createdAt.Format(time.RFC3339)
		thread.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.ThreadsSnapshot{
			Item: &api.ThreadsSnapshot_Thread{
				Thread: thread,
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

func (r *Dumper) Restore(ctx context.Context, snapshot *api.ThreadsSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.ThreadsSnapshot_Thread:

		query := "INSERT INTO %s(id, name, private, user_id, parent_id, next_id, prev_id, created_at, updated_at, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"

		_, err = r.pool.Exec(ctx, r.threadsTable(query), v.Thread.Id, v.Thread.Name, v.Thread.Private, v.Thread.UserId, v.Thread.ParentId,
			v.Thread.NextId, v.Thread.PrevId, v.Thread.CreatedAt, v.Thread.UpdatedAt, v.Thread.Description)

		return

	default:
		return errors.New("unknown snapshot item")
	}
}

func (r Dumper) threadsTable(query string) string {
	return fmt.Sprintf(query, r.threadsTableName)
}

