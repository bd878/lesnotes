package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
)

type FilesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewFilesRepository(tableName string, pool *pgxpool.Pool) *FilesRepository {
	return &FilesRepository{tableName: tableName, pool: pool}
}

func (r *FilesRepository) SaveMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[SaveMessageFiles]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	const insert = "INSERT INTO %s(file_id, message_id, user_id) VALUES ($1, $2, $3)"
	for _, fileID := range fileIDs {
		_, err = tx.Exec(ctx, r.table(insert), fileID, messageID, userID)
		if err != nil {
			return
		}
	}

	return
}

func (r *FilesRepository) ReadMessageFiles(ctx context.Context, messageID int64, userIDs []int64) (fileIDs []int64, err error) {
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
			fmt.Fprintf(os.Stderr, "[ReadMessageFiles]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	fileIDs = make([]int64, 0)

	ids := "$2"
	for i := 1; i < len(userIDs); i++ {
		ids += fmt.Sprintf(",$%d", i+2)
	}

	list := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		list[i] = id
	}

	var query = "SELECT file_id FROM %s WHERE message_id = $1 AND (user_id IN (" + ids + "))"

	rows, err := tx.Query(ctx, r.table(query), append([]interface{}{messageID}, list...)...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var fileID int64

		err = rows.Scan(&fileID)
		if err != nil {
			return
		}

		fileIDs = append(fileIDs, fileID)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *FilesRepository) UpdateMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[UpdateMessageFiles]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	const delete = "DELETE FROM %s WHERE message_id = $1 AND user_id = $2"

	_, err = tx.Exec(ctx, r.table(delete), messageID, userID)
	if err != nil {
		return
	}

	const insert = "INSERT INTO %s(file_id, message_id, user_id) VALUES ($1, $2, $3)"
	for _, fileID := range fileIDs {
		_, err = tx.Exec(ctx, r.table(insert), fileID, messageID, userID)
		if err != nil {
			return
		}
	}

	return
}

func (r *FilesRepository) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	const delete = "DELETE FROM %s WHERE message_id = $1 AND user_id = $2"

	_, err = r.pool.Exec(ctx, r.table(delete), id, userID)

	return
}

func (r *FilesRepository) DeleteMessage(ctx context.Context, messageID, userID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[DeleteMessage]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	const delete = "DELETE FROM %s WHERE message_id = $1 AND user_id = $2"
	_, err = tx.Exec(ctx, r.table(delete), messageID, userID)

	return
}

func (r *FilesRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping files repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *FilesRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring files repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r FilesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}