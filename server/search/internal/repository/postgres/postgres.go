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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.messagesTable(query), id, userID, name, title, text, private)

	return nil
}

func (r *Repository) UpdateMessage(ctx context.Context, id, userID int64, name, title, text string, private int32) (err error) {
	return nil
}

func (r *Repository) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	return nil
}

func (r *Repository) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	return nil
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
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, r.messagesTable(query), id, userID)

	return nil
}

func (r *Repository) SearchMessages(ctx context.Context, userID int64, substr string) (list []*searchmodel.Message, err error) {
	const query = "SELECT id, name, title, text, private FROM %s WHERE user_id = $1 AND name || ' ' || title || ' ' || text @@ $2"

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

	var rows pgx.Rows

	rows, err = tx.Query(ctx, r.messagesTable(query), userID, substr)
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