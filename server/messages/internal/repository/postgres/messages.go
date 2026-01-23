package postgres

import (
	"io"
	"fmt"
	"os"
	"strings"
	"time"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type MessagesRepository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func NewMessagesRepository(tableName string, pool *pgxpool.Pool) *MessagesRepository {
	return &MessagesRepository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *MessagesRepository) Create(ctx context.Context, id int64, text string, title string, userID int64, private bool, name string) (err error) {
	const query = "INSERT INTO %s(id, text, private, name, user_id, title) VALUES ($1, $2, $3, $5, $6, $7)"

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
			fmt.Fprintf(os.Stderr, "[Create]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table(query), id, text, private, name, userID, title)

	return
}

func (r *MessagesRepository) Update(ctx context.Context, userID, id int64, newText, newTitle, newName string) (err error) {
	const query = "UPDATE %s SET text = $3, title = $4, name = $5 WHERE user_id = $1 AND id = $2"
	const selectQuery = "SELECT text, title, name FROM %s WHERE user_id = $1 AND id = $2"

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
			fmt.Fprintf(os.Stderr, "[Update]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var text, title, name  string

	err = tx.QueryRow(ctx, r.table(selectQuery), userID, id).Scan(&text, &title, &name)
	if err != nil {
		return
	}

	if newText != "" {
		text = newText
	}

	if newTitle != "" {
		title = newTitle
	}

	if newName != "" {
		name = newName
	}

	_, err = tx.Exec(ctx, r.table(query), userID, id, text, title, name)

	return
}

func (r *MessagesRepository) DeleteMessage(ctx context.Context, userID, id int64) (err error) {
	_, err = r.pool.Exec(ctx, r.table("DELETE FROM %s WHERE id = $1 AND user_id = $2"), id, userID)

	return
}

func (r *MessagesRepository) Publish(ctx context.Context, userID int64, ids []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[Publish]: rollback with error: %v\n", err)
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

func (r *MessagesRepository) Private(ctx context.Context, userID int64, ids []int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[Private]: rollback with error: %v\n", err)
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

func (r *MessagesRepository) Read(ctx context.Context, userIDs []int64, id int64, name string) (message *model.Message, err error) {
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
			fmt.Fprintf(os.Stderr, "[Read]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	message = &model.Message{}

	var createdAt, updatedAt time.Time

	ids := "$2"
	for i := 1; i < len(userIDs); i++ {
		ids += fmt.Sprintf(",$%d", i+2)
	}

	list := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		list[i] = id
	}

	if id != 0 {

		err = tx.QueryRow(ctx, r.table(`
SELECT id, user_id, created_at, updated_at, text, private, name, title FROM %s WHERE id = $1 AND (user_id IN (` + ids + `) OR private = false)
`), append([]interface{}{id}, list...)...).Scan(&message.ID, &message.UserID, &createdAt, &updatedAt, &message.Text, &message.Private, &message.Name, &message.Title)

	} else if name != "" {

		err = tx.QueryRow(ctx, r.table(`
SELECT id, user_id, created_at, updated_at, text, private, name, title FROM %s WHERE name = $1 AND (user_id IN (` + ids + `) OR private = false)
`), append([]interface{}{name}, list...)...).Scan(&message.ID, &message.UserID, &createdAt, &updatedAt, &message.Text, &message.Private, &message.Name, &message.Title)

	}

	if err != nil {
		return
	}

	message.CreateUTCNano = createdAt.UnixNano()
	message.UpdateUTCNano = updatedAt.UnixNano()

	return
}

func (r *MessagesRepository) ReadBatchMessages(ctx context.Context, userID int64, messageIDs []int64) (messages []*model.Message, err error) {
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
			fmt.Fprintf(os.Stderr, "[ReadBatchMessages]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	messages = make([]*model.Message, 0)
	if len(messageIDs) == 0 {
		return
	}

	var rows pgx.Rows

	var order string
	pairs := make([]string, len(messageIDs))
	for i, messageID := range messageIDs {
		pairs[i] = fmt.Sprintf("(%d, %d)", messageID, i)
	}
	order = strings.Join(pairs, ",")

	query := r.table(`SELECT m.id, m.user_id, m.name, m.text, m.private, m.created_at, m.updated_at, m.title FROM %s m `) +
		fmt.Sprintf(` JOIN (VALUES %s) AS x(id, ordering) ON m.id = x.id WHERE m.user_id = $1 ORDER BY x.ordering ASC`, order)

	rows, err = tx.Query(ctx, query, userID)
	if err != nil {
		return
	}

	for rows.Next() {
		message := &model.Message{}

		var createdAt, updatedAt time.Time

		err = rows.Scan(&message.ID, &message.UserID, &message.Name, &message.Text, &message.Private, &createdAt, &updatedAt, &message.Title)
		if err != nil {
			return
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

func (r *MessagesRepository) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[DeleteUserMessages]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table("DELETE FROM %s WHERE user_id = $1"), userID)
	return
}

func (r *MessagesRepository) ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error) {
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
			fmt.Fprintf(os.Stderr, "[ReadMessages]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows

	query := "SELECT id, user_id, name, text, private, created_at, updated_at, title FROM %s WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err = tx.Query(ctx, r.table(query), userID, limit, offset)
	if err != nil {
		return
	}

	defer rows.Close()

	messages = make([]*model.Message, 0)
	for rows.Next() {
		message := &model.Message{}

		var createdAt, updatedAt time.Time

		err = rows.Scan(&message.ID, &message.UserID, &message.Name, &message.Text, &message.Private, &createdAt, &updatedAt, &message.Title)
		if err != nil {
			return
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
		err = tx.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE user_id = $1"), userID).Scan(&count)
		if err != nil {
			return
		}

		if count <= offset + limit {
			isLastPage = true
		}
	}

	return
}

func (r *MessagesRepository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.table("TRUNCATE TABLE %s"))
	return
}

func (r *MessagesRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping messages repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *MessagesRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring messages repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r MessagesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}