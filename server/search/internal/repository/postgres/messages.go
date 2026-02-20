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

type MessagesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewMessagesRepository(pool *pgxpool.Pool, tableName string) *MessagesRepository {
	return &MessagesRepository{tableName: tableName, pool: pool}
}

func (r *MessagesRepository) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
	const query = "INSERT INTO %s(id, user_id, name, title, text, private) VALUES ($1, $2, $3, $4, $5, $6)"

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
			fmt.Fprintf(os.Stderr, "[SaveMessage]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.table(query), id, userID, name, title, text, private)

	return
}

func (r *MessagesRepository) UpdateMessage(ctx context.Context, id, userID int64, newName, newTitle, newText string) (err error) {
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
			fmt.Fprintf(os.Stderr, "[UpdateMessage]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var text, title, name string

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

func (r *MessagesRepository) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[PublishMessages]: rollback with error: %v\n", err)
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

func (r *MessagesRepository) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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
			fmt.Fprintf(os.Stderr, "[PrivateMessages]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	for _, id := range ids {
		_, err = tx.Exec(ctx, r.table("UPDATE %s SET private = true WHERE user_id = $1 AND id = $2"), userID, id)
		if err != nil {
			return
		}
	}

	return
}

func (r *MessagesRepository) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1 AND user_id = $2"

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

	_, err = tx.Exec(ctx, r.table(query), id, userID)

	return
}

func (r *MessagesRepository) SearchMessages(ctx context.Context, userID int64, substr string, public int) (list []*search.Message, err error) {
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
			fmt.Fprintf(os.Stderr, "[SearchMessages]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows

	if public != -1 {
		var private bool
		if public == 0 {
			private = true
		} else {
			private = false
		}

		rows, err = tx.Query(ctx, r.table("SELECT id, name, title, text, private FROM %s WHERE user_id = $1 AND private = $2 AND name || ' ' || title || ' ' || text ILIKE $3"), userID, private, "'%" + substr + "%'")
	} else {
		rows, err = tx.Query(ctx, r.table("SELECT id, name, title, text, private FROM %s WHERE user_id = $1 AND name || ' ' || title || ' ' || text ILIKE $2"), userID, "'%" + substr + "%'")
	}

	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*search.Message, 0)
	for rows.Next() {
		message := &search.Message{
			UserID: userID,
		}

		err = rows.Scan(&message.ID, &message.Name, &message.Title, &message.Text, &message.Private)
		if err != nil {
			return
		}

		list = append(list, message)
	}

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