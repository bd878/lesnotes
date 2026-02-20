package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	search "github.com/bd878/gallery/server/search/pkg/model"
)

type ThreadsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewThreadsRepository(pool *pgxpool.Pool, tableName string) *ThreadsRepository {
	return &ThreadsRepository{tableName: tableName, pool: pool}
}

func (r *ThreadsRepository) SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool) (err error) {
	const query = "INSERT INTO %s(id, user_id, parent_id, name, description, private) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID, parentID, name, description, private)

	return
}

func (r *ThreadsRepository) SearchThreads(ctx context.Context, parentID, userID int64) (list []*search.Thread, err error) {
	const query = "SELECT id, user_id, description, private FROM %s WHERE user_id = $1 AND parent_id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[SearchThreads]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows

	rows, err = r.pool.Query(ctx, r.table(query), userID, parentID)
	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*search.Thread, 0)
	for rows.Next() {
		thread := &search.Thread{
			UserID: userID,
		}

		err = rows.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.Private)
		if err != nil {
			return
		}

		list = append(list, thread)
	}

	return
}

func (r *ThreadsRepository) UpdateThread(ctx context.Context, id, userID int64, newName, newDescription string) (err error) {
	const query = "SELECT name, description FROM %s WHERE user_id = $1 AND id = $2"
	const update = "UPDATE %s SET name = $3, description = $4 WHERE user_id = $1 AND id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[UpdateThread]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var name, description string

	err = tx.QueryRow(ctx, r.table(query), userID, id).Scan(&name, &description)
	if err != nil {
		return
	}

	if newName != "" {
		name = newName
	}

	if newDescription != "" {
		description = newDescription
	}

	_, err = tx.Exec(ctx, r.table(update), userID, id, name, description)

	return
}

func (r *ThreadsRepository) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(query), userID, id)

	return
}

func (r *ThreadsRepository) ChangeThreadParent(ctx context.Context, id, userID, newParentID int64) (err error) {
	const query = "SELECT parent_id FROM %s WHERE user_id = $1 AND id = $2"
	const update = "UPDATE %s SET parent_id = $3 WHERE user_id = $1 AND id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[ChangeThreadParent]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var parentID int64

	if newParentID != 0 {
		parentID = newParentID
	}

	err = tx.QueryRow(ctx, r.table(query), userID, id).Scan(&parentID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(update), userID, id, parentID)

	return
}

func (r *ThreadsRepository) PublishThread(ctx context.Context, id, userID int64) (err error) {
	const update = "UPDATE %s SET private = false WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *ThreadsRepository) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	const update = "UPDATE %s SET private = true WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *ThreadsRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping threads repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *ThreadsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring threads repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r ThreadsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}