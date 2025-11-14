package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
)

type Repository struct {
	tableName          string
	pool               *pgxpool.Pool
}

func New(pool *pgxpool.Pool, tableName string) *Repository {
	return &Repository{tableName, pool}
}


func (r *Repository) CreateThread(ctx context.Context, id, userID int64, parentID int64, name string, private bool) (err error) {
	const query = "INSERT INTO %s(id, user_id, parent_id, name, private) VALUES ($1, $2, $3, $4, $5)"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[CreateThread]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table(query), id, userID, parentID, name, private)

	return nil
}

func (r *Repository) UpdateThread(ctx context.Context, id, userID, newParentID int64, newName string, newPrivate int32) (err error) {
	const query = "UPDATE %s SET parent_id = $3, private = $4, name = $5 WHERE user_id = $1 AND id = $2"
	const selectQuery = "SELECT parent_id, private, name FROM %s WHERE user_id = $1 AND id = $2"

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

	var (
		name  string
		parentID int64
		private  bool
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), userID, id).Scan(&parentID, &private, &name)
	if err != nil {
		return
	}

	if newName != "" {
		name = newName
	}

	if newParentID != -1 {
		parentID = newParentID
	}

	if newPrivate != -1 {
		if newPrivate == 0 {
			private = false
		} else if newPrivate == 1 {
			private = true
		}
	}

	_, err = tx.Exec(ctx, r.table(query), userID, parentID, private, name)
	if err != nil {
		return
	}

	return
}

func (r *Repository) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[PrivateThread]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true WHERE user_id = $1 AND id = $2"), userID, id)

	return
}

func (r *Repository) PublishThread(ctx context.Context, id, userID int64) (err error) {
	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[PublishThread]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table("UPDATE %s SET private = false WHERE user_id = $1 AND id = $2"), userID, id)

	return
}

func (r *Repository) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[DeleteThread]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var parentID int64
	err = tx.QueryRow(ctx, r.table("SELECT parent_id FROM %s WHERE id = $1 AND user_id = $2"), id, userID).Scan(&parentID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("UPDATE %s SET parent_id = $3 WHERE user_id = $1 AND parent_id = $2"), userID, id, parentID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE id = $1 AND user_id = $2"), id, userID)

	return
}

func (r *Repository) ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error) {
	const query = "SELECT parent_id FROM %s WHERE user_id = $1 AND id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[ResolvePath]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	threadID := id

	ids = make([]int64, 0)

	for threadID != 0 {
		var parentID int64

		err = tx.QueryRow(ctx, r.table(query), userID, threadID).Scan(&parentID)
		if err != nil {
			return
		}

		ids = append(ids, threadID)

		threadID = parentID
	}

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