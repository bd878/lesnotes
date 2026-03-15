package postgres

import (
	"os"
	"time"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type CommentsRepository struct {
	tableName    string
	pool         *pgxpool.Pool
}

func NewCommentsRepository(pool *pgxpool.Pool, tableName string) *CommentsRepository {
	return &CommentsRepository{
		tableName:     tableName,
		pool:          pool,
	}
}

func (r *CommentsRepository) Create(ctx context.Context, id, userID, messageID int64, text string, metadata []byte, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, message_id, user_id, text, metadata, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"

	_, err = r.pool.Exec(ctx, r.table(query), id, messageID, userID, text, metadata, createdAt, updatedAt)

	return
}

func (r *CommentsRepository) Update(ctx context.Context, id, userID int64, text *string, updatedAt string) (err error) {
	const query = "UPDATE %s SET text = $3, updated_at = $4 WHERE id = $1 AND user_id = $2"

	logger.Debugw("update", "id", id, "user_id", userID, "text", *text, "updated_at", updatedAt)

	_, err = r.pool.Exec(ctx, r.table(query), id, userID, text, updatedAt)

	return
}

func (r *CommentsRepository) Delete(ctx context.Context, id, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1 AND user_id = $2"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID)

	return
}

func (r *CommentsRepository) DeleteMessageComments(ctx context.Context, messageID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[DeleteMessageComments]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	const query = "DELETE FROM %s WHERE message_id = $1"

	_, err = tx.Exec(ctx, r.table(query), messageID)

	return
}

func (r *CommentsRepository) Read(ctx context.Context, id, userID int64) (comment *api.Comment, err error) {
	const query = "SELECT message_id, text, metadata, created_at, updated_at FROM %s WHERE id = $1 AND user_id = $2"

	comment = &api.Comment{
		Id:       id,
		UserId:   userID,
	}

	var createdAt, updatedAt time.Time
	err = r.pool.QueryRow(ctx, r.table(query), id, userID).Scan(&comment.MessageId, &comment.Text, &comment.Metadata,
		&createdAt, &updatedAt)
	if err != nil {
		return
	}

	comment.CreatedAt = createdAt.Format(time.RFC3339)
	comment.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *CommentsRepository) ListMessageComments(ctx context.Context, messageID int64, limit, offset int32) (list *api.CommentsList, err error) {
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
			fmt.Fprintf(os.Stderr, "[ListMessageComments]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows
	const query = "SELECT id, user_id, text, metadata, created_at, updated_at FROM %s WHERE message_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err = tx.Query(ctx, r.table(query), messageID, limit, offset)
	defer rows.Close()
	if err != nil {
		return
	}

	comments := make([]*api.Comment, 0)
	for rows.Next() {
		comment := &api.Comment{
			MessageId:    messageID,
		}

		var createdAt, updatedAt time.Time
		err = rows.Scan(&comment.Id, &comment.UserId, &comment.Text, &comment.Metadata, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		comment.CreatedAt = createdAt.Format(time.RFC3339)
		comment.UpdatedAt = updatedAt.Format(time.RFC3339)

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return
	}

	var (
		total int32
		isLastPage bool
	)

	err = tx.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE message_id = $1"), messageID).Scan(&total)
	if err != nil {
		return
	}

	if total <= offset + limit {
		isLastPage = true
	}

	list = &api.CommentsList{
		Comments:     comments,
		IsLastPage:   isLastPage,
		IsFirstPage:  int32(offset) == 0,
		Count:        int32(len(comments)),
		Total:        total,
	}

	return
}

func (r CommentsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}