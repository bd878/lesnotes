package repository

import (
	"fmt"
	"os"
	"time"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(tableName string, pool *pgxpool.Pool) *Repository {
	return &Repository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *Repository) Create(ctx context.Context, message *model.Message) (err error) {
	const query = "INSERT INTO %s(id, text, file_ids, private, name, user_id, thread_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"

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
			fmt.Fprintf(os.Stderr, "rollback with error: %v", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	fileIDs, err := json.Marshal(message.FileIDs)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, r.table(query), message.ID, message.Text, fileIDs, message.Private, message.Name, message.UserID, message.ThreadID)
	if err != nil {
		return
	}

	return nil
}

/**
 * newText == "" : left as is
 * newThreadID == -1 : left as is
 * newPrivate == -1 : left as is
 * @param  {[type]} r *Repository)  Update(ctx context.Context, id int32, userID int32, text string, threadID int32, fileIDs []int32, private int) (*model.UpdateMessageResult, error [description]
 * @return {error}   error
 */
func (r *Repository) Update(ctx context.Context, id int32, userID int32, newText string, newThreadID int32, newFileIDs []int32, newPrivate int) (err error) {
	const query = "UPDATE %s SET text = $3, thread_id = $4, file_ids = $5, private = $6 WHERE id = $1 AND user_id = $2"
	const selectQuery = "SELECT text, thread_id, file_ids, private FROM %s WHERE id = $1 AND user_id = $2"

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
			fmt.Fprintf(os.Stderr, "rollback with error: %v", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var (
		text     string
		threadID int32
		fileIDs  []byte
		private  bool
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), id, userID).Scan(&text, &threadID, &fileIDs, &private)
	if err != nil {
		return
	}

	if newText != "" {
		text = newText
	}

	if newThreadID != -1 {
		threadID = newThreadID
	}

	if newFileIDs != nil {
		fileIDs, err = json.Marshal(newFileIDs)
		if err != nil {
			return
		}
	}

	if newPrivate != -1 {
		if newPrivate == 0 {
			private = false
		} else if newPrivate == 1 {
			private = true
		}
	}

	_, err = tx.Exec(ctx, r.table(query), id, userID, text, threadID, fileIDs, private)
	if err != nil {
		return
	}

	return nil
}

/**
 * Delete message and move ancestor messages on current thread
 * @param  {[type]} r *Repository)  DeleteMessage(ctx context.Context, userID, id, parentThreadID int32) (err error [description]
 * @return {error}   error
 */
func (r *Repository) DeleteMessage(ctx context.Context, userID, id int32) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var threadID int32
	err = tx.QueryRow(ctx, r.table("SELECT thread_id FROM %s WHERE id = $1 AND user_id = $2"), id, userID).Scan(&threadID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("UPDATE %s SET thread_id = $3 WHERE user_id = $1 AND thread_id = $2"), userID, id, threadID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE id = $1 AND user_id = $2"), id, userID)
	if err != nil {
		return
	}

	return nil
}

func (r *Repository) Publish(ctx context.Context, userID, id int32) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = false WHERE id = $1 AND user_id = $2"), id, userID)
	return
}

func (r *Repository) Private(ctx context.Context, userID, id int32) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true WHERE id = $1 AND user_id = $2"), id, userID)
	return
}

func (r *Repository) Read(ctx context.Context, userIDs []int32, id int32) (message *model.Message, err error) {
	message = &model.Message{ID: id}

	var (
		fileIDs   []byte
		createdAt, updatedAt time.Time
	)

	ids := "$2"
	for i := 1; i < len(userIDs); i++ {
		ids += fmt.Sprintf(",%d", i+2)
	}

	list := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		list[i] = id
	}

	err = r.pool.QueryRow(ctx, r.table(`
SELECT user_id, thread_id, file_ids, created_at, updated_at, text, private, name FROM %s WHERE id = $1 AND (user_id IN (` + ids + `) OR private = 0)
`), append([]interface{}{id}, list...)...).Scan(&message.UserID, &message.ThreadID, &fileIDs, &createdAt, &updatedAt, &message.Text, &message.Private, &message.Name)
	if err != nil {
		return
	}

	if fileIDs != nil {
		err = json.Unmarshal(fileIDs, &message.FileIDs)
		if err != nil {
			return
		}
	}

	message.CreateUTCNano = createdAt.UnixNano()
	message.UpdateUTCNano = updatedAt.UnixNano()

	return
}

func (r *Repository) DeleteAllUserMessages(ctx context.Context, userID int32) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE user_id = $1"), userID)
	return
}

func (r *Repository) ReadMessages(ctx context.Context, userID, threadID, limit, offset int32, private int32) (messages []*model.Message, isLastPage bool, err error) {
	var rows pgx.Rows

	query := "SELECT id, user_id, thread_id, file_ids, name, text, private, created_at, updated_at FROM %s WHERE user_id = $1 AND thread_id = $2 AND private = $5 ORDER BY created_at DESC LIMIT $3 OFFSET $4"

	rows, err = r.pool.Query(ctx, r.table(query), userID, threadID, limit, offset, private)
	defer rows.Close()

	messages = make([]*model.Message, 0)
	for rows.Next() {
		message := &model.Message{}

		var (
			fileIDs []byte
			createdAt, updatedAt time.Time
		)

		err = rows.Scan(&message.ID, &message.UserID, &message.ThreadID, &fileIDs, &message.Name, &message.Text, &message.Private, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		if fileIDs != nil {
			err = json.Unmarshal(fileIDs, &message.FileIDs)
			if err != nil {
				return
			}
		}

		message.CreateUTCNano = createdAt.UnixNano()
		message.UpdateUTCNano = updatedAt.UnixNano()

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return
	}

	if int32(len(messages)) < limit {
		isLastPage = true
	} else {
		var count int32
		err = r.pool.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE user_id = $1 AND thread_id = $2"), userID, threadID).Scan(&count)
		if err != nil {
			return
		}

		if count <= offset + limit {
			isLastPage = true
		}
	}

	return
}

func (r *Repository) Truncate(ctx context.Context) (err error) {
	_, err = r.pool.Exec(ctx, r.table("DELETE FROM %s"))
	return nil
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}