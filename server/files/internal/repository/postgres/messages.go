package postgres

import (
	"fmt"
	"os"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
)

type MessagesRepository struct {
	tableName    string
	pool         *pgxpool.Pool
}

func NewMessagesRepository(pool *pgxpool.Pool, tableName string) *MessagesRepository {
	return &MessagesRepository{
		tableName:     tableName,
		pool:          pool,
	}
}

func (r *MessagesRepository) SaveMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error) {
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

func (r *MessagesRepository) ReadMessageFiles(ctx context.Context, messageID int64, userIDs []int64) (fileIDs []int64, err error) {
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
	defer rows.Close()
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

func (r *MessagesRepository) DeleteFiles(ctx context.Context, ids []int64) (err error) {
	for _, id := range ids {
		_, err = r.pool.Exec(ctx, r.table("DELETE FROM %s WHERE file_id = $1"), id)
		if err != nil {
			logger.Errorln(err)
			continue
		}
	}

	return
}

func (r *MessagesRepository) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	_, err = r.pool.Exec(ctx, r.table("DELETE FROM %s WHERE message_id = $1 AND user_id = $2"), id, userID)

	return
}

func (r *MessagesRepository) UpdateMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error) {
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

func (r MessagesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
