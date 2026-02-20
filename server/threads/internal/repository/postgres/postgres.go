package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Repository struct {
	tableName          string
	pool               *pgxpool.Pool
}

func New(pool *pgxpool.Pool, tableName string) *Repository {
	return &Repository{tableName, pool}
}

func (r *Repository) ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error) {
	thread = &threads.Thread{}

	if id != 0 {
		const query = "SELECT id, parent_id, user_id, next_id, prev_id, name, description, private FROM %s WHERE id = $1 AND (user_id = $2 OR private = false)"

		err = r.pool.QueryRow(ctx, r.table(query), id, userID).Scan(&thread.ID, &thread.ParentID, &thread.UserID, &thread.NextID, &thread.PrevID,
			&thread.Name, &thread.Description, &thread.Private)

	} else if name != "" {
		const query = "SELECT id, parent_id, user_id, next_id, prev_id, name, description, private FROM %s WHERE name = $1 AND (user_id = $2 OR private = false)"

		err = r.pool.QueryRow(ctx, r.table(query), name, userID).Scan(&thread.ID, &thread.ParentID, &thread.UserID, &thread.NextID, &thread.PrevID,
			&thread.Name, &thread.Description, &thread.Private)
	}
	if err != nil {
		return
	}

	var total int32
	err = r.pool.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE user_id = $1 AND parent_id = $2"), thread.UserID, thread.ID).Scan(&total)
	if err != nil {
		return
	}

	thread.Count = total

	return
}

func (r *Repository) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*threads.Thread, isLastPage bool, err error) {
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
			tx.Rollback(ctx)
		default:
			tx.Commit(ctx)
		}
	}()

	const selectThreads = "SELECT id, name, private, next_id, prev_id FROM %s WHERE user_id = $1 AND parent_id = $2"

	var rows pgx.Rows
	rows, err = tx.Query(ctx, r.table(selectThreads), userID, parentID)
	if err != nil {
		return
	}

	unordered := make([]*threads.Thread, 0)
	for rows.Next() {
		thread := &threads.Thread{
			ParentID: parentID,
			UserID:   userID,
		}

		err = rows.Scan(&thread.ID, &thread.Name, &thread.Private, &thread.NextID, &thread.PrevID)
		if err != nil {
			return
		}

		unordered = append(unordered, thread)
	}

	var nextID int64
	list = make([]*threads.Thread, 0)
	for range unordered {
		for _, thread := range unordered {
			if thread.NextID == nextID {
				list = append(list, thread)
				nextID = thread.ID
				break
			}
		}
	}

	if err = rows.Err(); err != nil {
		return
	}

	for _, thread := range list {
		err = tx.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE user_id = $1 AND parent_id = $2"), userID, thread.ID).Scan(&thread.Count)
		if err != nil {
			return
		}
	}

	if int32(len(list)) <= offset+limit {
		isLastPage = true
	}

	end := min(int32(len(list)), offset+limit)

	list = list[offset:end]

	return
}

func (r *Repository) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	const query = "SELECT COUNT(*) FROM %s WHERE user_id = $1 AND parent_id = $2"

	err = r.pool.QueryRow(ctx, r.table(query), userID, id).Scan(&total)

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

	// validate arguments

	const updateNextThread = "UPDATE %s SET prev_id = $3 WHERE user_id = $1 AND id = $2"
	const updatePrevThread = "UPDATE %s SET next_id = $3 WHERE user_id = $1 AND id = $2"
	const selectParentID = "SELECT parent_id FROM %s WHERE user_id = $1 AND id = $2"
	const updateMe = "UPDATE %s SET parent_id = $3, next_id = $4, prev_id = $5 WHERE user_id = $1 AND id = $2"

	var currentParentID, currentNextID, currentPrevID, prevIDParent, nextIDParent int64
	err = tx.QueryRow(ctx, r.table("SELECT parent_id, next_id, prev_id FROM %s WHERE user_id = $1 AND id = $2"), userID, id).
		Scan(&currentParentID, &currentNextID, &currentPrevID)
	if err != nil {
		logger.Debugw("failed to select thread", "error", err)
		return
	}

	if prevID != 0 && nextID != 0 {
		return fmt.Errorf("either prevID or nextID must be given, not both; next_id = %d, prev_id = %d", prevID, nextID)
	}

	if prevID == id {
		return fmt.Errorf("prevID == id: %d, %d", prevID, id)
	}
	if nextID == id {
		return fmt.Errorf("nextID == id: %d, %d", nextID, id)
	}
	if parentID == id {
		return fmt.Errorf("parentID == id: %d, %d", parentID, id)
	}

	if prevID != 0 {
		err = tx.QueryRow(ctx, r.table(selectParentID), userID, prevID).Scan(&prevIDParent)
		if err != nil {
			logger.Debugw("failed to get prev id parent", "error", err)
			return
		}
	}

	if nextID != 0 {
		err = tx.QueryRow(ctx, r.table(selectParentID), userID, nextID).Scan(&nextIDParent)
		if err != nil {
			logger.Debugw("failed to get next id parent", "error", err)
			return
		}
	}

	if parentID != -1 {
		if prevID != 0 && prevIDParent != parentID {
			return fmt.Errorf("prevIDParent != parentID: %d != %d", prevIDParent, parentID)
		}

		if nextID != 0 && nextIDParent != parentID {
			return fmt.Errorf("nextIDParent != parentID: %d != %d", nextIDParent, parentID)
		}
	} else {
		if prevID != 0 {
			parentID = prevIDParent
		}

		if nextID != 0 {
			parentID = nextIDParent
		}
	}

	// validate parent
	threadID := parentID
	for threadID != 0 {
		if id == threadID {
			return fmt.Errorf("new parent %s is a relative of thread %s", parentID, id)
		}

		err = tx.QueryRow(ctx, r.table("SELECT parent_id FROM %s WHERE user_id = $1 AND id = $2"),
			userID, threadID).Scan(&threadID)
		if err != nil {
			return
		}
	}
	// end validate parent

	// end validate arguments
	// unlink

	if currentPrevID != 0 {
		_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, currentPrevID, currentNextID)
		if err != nil {
			logger.Debugw("failed to update prev thread", "error", err)
			return
		}
	}

	if currentNextID != 0 {
		_, err = tx.Exec(ctx, r.table(updateNextThread), userID, currentNextID, currentPrevID)
		if err != nil {
			logger.Debugw("failed to update next thread", "error", err)
			return
		}
	}

	// end unlink
	// reorder
	if prevID != 0 {
		// reorder before

		var nextID int64

		err = tx.QueryRow(ctx, r.table("SELECT next_id FROM %s WHERE user_id = $1 AND id = $2"), userID, prevID).
			Scan(&nextID)
		if err != nil {
			logger.Debugw("failed to get next id of prev id", "prev_id", prevID, "error", err)
			return
		}

		_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, prevID, id)
		if err != nil {
			return
		}

		_, err = tx.Exec(ctx, r.table(updateNextThread), userID, nextID, id)
		if err != nil {
			return
		}

		_, err = tx.Exec(ctx, r.table(updateMe), userID, id, parentID, nextID, prevID)

		return

	} else if nextID != 0 {
		// reorder after

		var prevID int64

		err = tx.QueryRow(ctx, r.table("SELECT prev_id FROM %s WHERE user_id = $1 AND id = $2"), userID, nextID).
			Scan(&prevID)
		if err != nil {
			logger.Debugw("failed to get prev id of next id", "next_id", nextID, "error", err)
			return
		}

		_, err = tx.Exec(ctx, r.table(updatePrevThread), userID, prevID, id)
		if err != nil {
			return
		}

		_, err = tx.Exec(ctx, r.table(updateNextThread), userID, nextID, id)
		if err != nil {
			return
		}

		_, err = tx.Exec(ctx, r.table(updateMe), userID, id, parentID, nextID, prevID)

		return

	}
	// end reorder

	return
}

func (r *Repository) AppendThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error) {
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

	const insert = "INSERT INTO %s(id, user_id, parent_id, name, description, private, next_id, prev_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	const selectLastThread = "SELECT id FROM %s WHERE user_id = $1 AND parent_id = $2 AND next_id = 0"
	const updateLastThread = "UPDATE %s SET next_id = $4 WHERE user_id = $1 AND id = $2 AND parent_id = $3"

	var lastThreadID int64
	err = tx.QueryRow(ctx, r.table(selectLastThread), userID, parentID).Scan(&lastThreadID)
	if err != nil && err != pgx.ErrNoRows {
		return
	}

	if err == pgx.ErrNoRows {
		// new thread
		_, err = tx.Exec(ctx, r.table(insert), id, userID, parentID, name, description, private, 0, 0)
		return
	}

	_, err = tx.Exec(ctx, r.table(updateLastThread), userID, lastThreadID, parentID, id)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(insert), id, userID, parentID, name, description, private, 0 /* next_id */, lastThreadID)

	return
}

func (r *Repository) UpdateThread(ctx context.Context, id, userID int64, newName, newDescription string) (err error) {
	const query = "UPDATE %s SET description = $3, name = $4 WHERE user_id = $1 AND id = $2"
	const selectQuery = "SELECT description, name FROM %s WHERE user_id = $1 AND id = $2"

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

	err = tx.QueryRow(ctx, r.table(selectQuery), userID, id).Scan(&description, &name)
	if err != nil {
		return
	}

	if newName != "" {
		name = newName
	}

	if newDescription != "" {
		description = newDescription
	}

	_, err = tx.Exec(ctx, r.table(query), userID, id, description, name)

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