package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

type Repository struct {
	messagesTableName  string
	filesTableName     string
	pool               *pgxpool.Pool
}

// TODO: split repo on messages and files, like billing invoices and payments
func New(pool *pgxpool.Pool, messagesTableName, filesTableName string) *Repository {
	return &Repository{messagesTableName, filesTableName, pool}
}

func (r *Repository) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
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

	_, err = tx.Exec(ctx, r.messagesTable(query), id, userID, name, title, text, private)

	return nil
}

func (r *Repository) UpdateMessage(ctx context.Context, id, userID int64, newName, newTitle, newText string) (err error) {
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

	err = tx.QueryRow(ctx, r.messagesTable(selectQuery), userID, id).Scan(&text, &title, &name)
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

	_, err = tx.Exec(ctx, r.messagesTable(query), userID, id, text, title, name)
	if err != nil {
		return
	}

	return
}

func (r *Repository) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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
		_, err = tx.Exec(ctx, r.messagesTable("UPDATE %s SET private = false WHERE user_id = $1 AND id = $2"), userID, id)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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
		_, err = r.pool.Exec(ctx, r.messagesTable("UPDATE %s SET private = true WHERE user_id = $1 AND id = $2"), userID, id)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
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

	_, err = tx.Exec(ctx, r.messagesTable(query), id, userID)

	return nil
}

func (r *Repository) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*searchmodel.Message, err error) {
	// TODO: filter by thread_id

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

		rows, err = tx.Query(ctx, r.messagesTable("SELECT id, name, title, text, private FROM %s WHERE user_id = $1 AND private = $2 AND name || ' ' || title || ' ' || text ILIKE $3"), userID, private, "%" + substr + "%")
	} else {
		rows, err = tx.Query(ctx, r.messagesTable("SELECT id, name, title, text, private FROM %s WHERE user_id = $1 AND name || ' ' || title || ' ' || text ILIKE $2"), userID, "%" + substr + "%")
	}

	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*searchmodel.Message, 0)
	for rows.Next() {
		message := &searchmodel.Message{
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

func (r *Repository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.messagesTable("TRUNCATE TABLE %s"))
	return
}

func (r *Repository) Dump(ctx context.Context) (reader io.ReadCloser, err error) {
	var (
		writer io.WriteCloser
		conn   *pgxpool.Conn
	)

	query := r.messagesTable("COPY %s TO STDOUT BINARY")

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

	query := r.messagesTable("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	if err != nil {
		conn.Release()
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)
	defer conn.Release()

	return
}

func (r Repository) messagesTable(query string) string {
	return fmt.Sprintf(query, r.messagesTableName)
}