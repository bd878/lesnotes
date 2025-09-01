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

func (r *Repository) Create(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var files []byte
	if fileIDs != nil {
		files, err = json.Marshal(fileIDs)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, r.table(query), id, text, files, private, name, userID, threadID)

	return
}

/**
 * newText == "" : left as is
 * newThreadID == -1 : left as is
 * newPrivate == -1 : left as is
 * @param  {[type]} r *Repository)  Update(ctx context.Context, userID, id int64, text string, threadID int64, fileIDs []int64, private int) (error [description]
 * @return {error}   error
 */
func (r *Repository) Update(ctx context.Context, userID, id int64, newText string, newThreadID int64, newFileIDs []int64, newPrivate int) (err error) {
	const query = "UPDATE %s SET text = $3, thread_id = $4, file_ids = $5, private = $6 WHERE user_id = $1 AND id = $2"
	const selectQuery = "SELECT text, thread_id, file_ids, private FROM %s WHERE user_id = $1 AND id = $2"

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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var (
		text     string
		threadID int64
		fileIDs  []byte
		private  bool
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), userID, id).Scan(&text, &threadID, &fileIDs, &private)
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

	_, err = tx.Exec(ctx, r.table(query), userID, id, text, threadID, fileIDs, private)
	if err != nil {
		return
	}

	return
}

/**
 * Delete message and move ancestor messages on current thread
 * @param  {[type]} r *Repository)  DeleteMessage(ctx context.Context, userID, id, parentThreadID int64) (err error [description]
 * @return {error}   error
 */
func (r *Repository) DeleteMessage(ctx context.Context, userID, id int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var threadID int64
	err = tx.QueryRow(ctx, r.table("SELECT thread_id FROM %s WHERE id = $1 AND user_id = $2"), id, userID).Scan(&threadID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("UPDATE %s SET thread_id = $3 WHERE user_id = $1 AND thread_id = $2"), userID, id, threadID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE id = $1 AND user_id = $2"), id, userID)

	return
}

func (r *Repository) Publish(ctx context.Context, userID int64, ids []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	for _, id := range ids {
		_, err = tx.Exec(ctx, r.table("UPDATE %s SET private = false WHERE user_id = $1 AND id = $2"), userID, id)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) Private(ctx context.Context, userID int64, ids []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	for _, id := range ids {
		_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true WHERE user_id = $1 AND id = $2"), userID, id)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) Read(ctx context.Context, userIDs []int64, id int64) (message *model.Message, err error) {
	message = &model.Message{ID: id}

	var (
		fileIDs   []byte
		createdAt, updatedAt time.Time
	)

	ids := "$2"
	for i := 1; i < len(userIDs); i++ {
		ids += fmt.Sprintf(",$%d", i+2)
	}

	list := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		list[i] = id
	}

	err = r.pool.QueryRow(ctx, r.table(`
SELECT user_id, thread_id, file_ids, created_at, updated_at, text, private, name FROM %s WHERE id = $1 AND (user_id IN (` + ids + `) OR private = false)
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

func (r *Repository) ReadBatchMessages(ctx context.Context, userID int64, messageIDs []int64) (messages []*model.Message, err error) {
	var rows pgx.Rows

	ids := "$2"
	for i := 1; i < len(messageIDs); i++ {
		ids += fmt.Sprintf(",$%d", i+2)
	}

	list := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		list[i] = id
	}

	rows, err = r.pool.Query(ctx, r.table(`
SELECT id, user_id, thread_id, file_ids, name, text, private, created_at, updated_at FROM %s WHERE user_id = $1 AND (id IN (` + ids + `))
`), append([]interface{}{userID}, list...)...)
	if err != nil {
		return
	}

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

	return
}

func (r *Repository) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE user_id = $1"), userID)
	return
}

/**
 * @param  {[type]} r *Repository)  ReadThreadMessages(ctx context.Context, userID, threadID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error [description]
 * @return {[type]}   [description]
 */
func (r *Repository) ReadThreadMessages(ctx context.Context, userID, threadID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error) {
	var rows pgx.Rows

	query := "SELECT id, user_id, thread_id, file_ids, name, text, private, created_at, updated_at FROM %s WHERE user_id = $1 AND thread_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4"

	rows, err = r.pool.Query(ctx, r.table(query), userID, threadID, limit, offset)
	defer rows.Close()
	if err != nil {
		return
	}

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

/**
 * Read all user messages from all threads
 * @param  {[type]} r *Repository)  ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error [description]
 * @return {[type]}   [description]
 */
func (r *Repository) ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error) {
	var rows pgx.Rows

	query := "SELECT id, user_id, thread_id, file_ids, name, text, private, created_at, updated_at FROM %s WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err = r.pool.Query(ctx, r.table(query), userID, limit, offset)

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
		err = r.pool.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE user_id = $1"), userID).Scan(&count)
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
	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}