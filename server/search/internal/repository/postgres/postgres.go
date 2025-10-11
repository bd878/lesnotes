package postgres

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	messagesTableName  string
	filesTableName     string
	pool               *pgxpool.Pool
}

func New(pool *pgxpool.Pool, messagesTableName, filesTableName string) *Repository {
	return &Repository{messagesTableName, filesTableName, pool}
}

func (r *Repository) SaveMessage(ctx context.Context, id, userID int64, name, title, text string) (err error) {
	const query = "INSERT INTO %s(id, user_id, name, title, text) VALUES ($1, $2, $3, $4, $5)"

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

	_, err = tx.Exec(ctx, r.messagesTable(query), id, userID, name, title, text)

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

func (r Repository) messagesTable(query string) string {
	return fmt.Sprintf(query, r.messagesTableName)
}