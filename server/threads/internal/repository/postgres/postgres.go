package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Repository struct {
	tableName          string
	pool               *pgxpool.Pool
}

func New(pool *pgxpool.Pool, tableName string) *Repository {
	return &Repository{tableName, pool}
}

func (r *Repository) ReadThread(ctx context.Context, id, userID int64) (thread *threads.Thread, err error) {
	const query = "SELECT parent_id, next_id, prev_id, name, private FROM %s WHERE user_id = $1 AND id = $2"

	thread = &threads.Thread{
		ID: id,
		UserID: userID,
	}

	err = r.pool.QueryRow(ctx, r.table(query), userID, id).Scan(&thread.ParentID, &thread.NextID, &thread.PrevID, &thread.Name, &thread.Private)

	return
}

func (r *Repository) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []int64, isLastPage bool, err error) {
	return
}

func (r *Repository) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
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
			tx.Rollback(ctx)
		default:
			tx.Commit(ctx)
		}
	}()

	const selectThread = "SELECT parent_id, next_id, prev_id FROM %s WHERE user_id = $1 AND id = $2 AND parent_id = $3"
	const updateNextThread = "UPDATE %s SET prev_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"
	const updatePrevThread = "UPDATE %s SET next_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"
	const updateMe = "UPDATE %s SET parent_id = $4, next_id = $5, prev_id = $6 WHERE user_id = $1 AND id = $2 AND parent_id = $3"

	// Unlink
	var currentParentID, currentNextID, currentPrevID int64
	err = tx.QueryRow(ctx, r.table(selectThread), userID, id, parentID).Scan(&currentParentID, &currentNextID, &currentPrevID)
	if err != nil {
		return
	}

	if currentPrevID != 0 {
		_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, id, parentID, currentNextID)
		if err != nil {
			return
		}
	}

	if currentNextID != 0 {
		_, err = tx.Exec(ctx, r.table(updateNextThread), userID, id, parentID, currentPrevID)
		if err != nil {
			return
		}
	}

	// End unlink
	// Link

	_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, id, id)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(updateNextThread), userID, id, id)
	if err != nil {
		return
	}

	// Update me
	if parentID == -1 { // 0 is root
		parentID = currentParentID
	}

	_, err = tx.Exec(ctx, r.table(updateMe), userID, id, currentParentID, parentID, nextID, prevID)

	return
}

func (r *Repository) AppendThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name string, private bool) (err error) {
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

	const insert = "INSERT INTO %s(id, user_id, parent_id, name, private, next_id, prev_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	const selectLastThread = "SELECT id FROM %s WHERE user_id = $1 AND parent_id = $2 AND next_id = 0"
	const updateLastThread = "UPDATE %s SET next_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"

	var lastThreadID int64
	row := tx.QueryRow(ctx, r.table(selectLastThread), userID, parentID)
	err = row.Scan(&lastThreadID)
	if err != nil && err != pgx.ErrNoRows {
		return
	}

	if err == pgx.ErrNoRows {
		// new thread
		_, err = tx.Exec(ctx, r.table(insert), id, userID, parentID, name, private, 0, 0)
		return
	}

	_, err = tx.Exec(ctx, r.table(updateLastThread), userID, lastThreadID, parentID, id)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(insert), id, userID, parentID, name, private, 0 /* next_id */, lastThreadID)

	return
}

func (r *Repository) UpdateThread(ctx context.Context, id, userID int64, newName string, newPrivate int32) (err error) {
	const query = "UPDATE %s SET private = $3, name = $4 WHERE user_id = $1 AND id = $2"
	const selectQuery = "SELECT private, name FROM %s WHERE user_id = $1 AND id = $2"

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
		private  bool
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), userID, id).Scan(&private, &name)
	if err != nil {
		return
	}

	if newName != "" {
		name = newName
	}

	if newPrivate != -1 {
		if newPrivate == 0 {
			private = false
		} else if newPrivate == 1 {
			private = true
		}
	}

	_, err = tx.Exec(ctx, r.table(query), userID, private, name)
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

	// Unlink
	const selectThread = "SELECT parent_id, next_id, prev_id FROM %s WHERE user_id = $1 AND id = $2"
	const updateNextThread = "UPDATE %s SET prev_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"
	const updatePrevThread = "UPDATE %s SET next_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"

	var parentID, nextID, prevID int64
	err = tx.QueryRow(ctx, r.table(selectThread), userID, id).Scan(&parentID, &nextID, &prevID)
	if err != nil {
		logger.Debugln("cannot select thread")
		return
	}

	if nextID != 0 {
		_, err = tx.Exec(ctx, r.table(updateNextThread), userID, nextID, parentID, prevID)
		if err != nil {
			logger.Debugln("cannot update next thread")
			return
		}
	}

	if prevID != 0 {
		_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, prevID, parentID, nextID)
		if err != nil {
			logger.Debugln("cannot update prev thread")
			return
		}
	}

	// End unlink
	// Delete
	const deleteThread = "DELETE FROM %s WHERE user_id = $1 AND id = $2"

	_, err = tx.Exec(ctx, r.table(deleteThread), userID, id)
	if err != nil {
		logger.Debugln("cannot delete thread")
		return
	}

	// End delete
	// Link children
	const countChildren = "SELECT COUNT(*) FROM %s WHERE user_id = $1 AND parent_id = $2"
	const selectLastParentThread = "SELECT id FROM %s WHERE user_id = $1 AND parent_id = $2 AND next_id = 0"
	const selectFirstThread = "SELECT id FROM %s WHERE user_id = $1 AND parent_id = $2 AND prev_id = 0"
	const updateFirstThread = "UPDATE %s SET prev_id = $3 WHERE user_id = $1 AND parent_id = $2 AND prev_id = 0"
	const updateLastParentThread = "UPDATE %s SET next_id = $3 WHERE user_id = $1 AND parent_id = $2 AND next_id = 0"
	const moveChildren = "UPDATE %s SET parent_id = $3 WHERE user_id = $1 AND parent_id = $2"

	var count int32
	err = tx.QueryRow(ctx, r.table(countChildren), userID, id).Scan(&count)
	if err != nil {
		logger.Debugln("cannot count children")
		if err == pgx.ErrNoRows {
			return nil
		}
		return
	}

	// nothing to move, exit
	if count == 0 {
		return nil
	}

	var lastParentThreadID, firstThreadID int64
	err = tx.QueryRow(ctx, r.table(selectLastParentThread), userID, parentID).Scan(&lastParentThreadID)
	if err != nil && err != pgx.ErrNoRows {
		logger.Debugln("cannot select last parent thread")
		return
	}

	// it must have children, since count > 0
	err = tx.QueryRow(ctx, r.table(selectFirstThread), userID, id).Scan(&firstThreadID)
	if err != nil {
		logger.Debugln("cannot select first thread")
		return
	}

	_, err = tx.Exec(ctx, r.table(updateFirstThread), userID, id, lastParentThreadID)
	if err != nil {
		logger.Debugln("cannot update first thread")
		return
	}

	if lastParentThreadID != 0 {
		// move to non-empty thread
		_, err = tx.Exec(ctx, r.table(updateLastParentThread), userID, parentID, firstThreadID)
		if err != nil {
			logger.Debugln("cannot update last parent thread")
			return
		}
	}

	_, err = tx.Exec(ctx, r.table(moveChildren), userID, id, parentID)
	if err != nil {
		logger.Debugln("cannot move children")
		return
	}

	// End link children

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