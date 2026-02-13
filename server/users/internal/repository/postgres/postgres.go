package repository

import (
	"fmt"
	"io"
	"os"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(pool *pgxpool.Pool, tableName string) *Repository {
	return &Repository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *Repository) Save(ctx context.Context, id int64, login, salt string, metadata []byte) (err error) {
	const query = "INSERT INTO %s(id, login, salt, metadata) VALUES ($1, $2, $3, $4)"

	_, err = r.pool.Exec(ctx, r.table(query), id, login, salt, metadata)

	return
}

func (r *Repository) Delete(ctx context.Context, id int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), id)

	return
}

/**
 * Find by id or login
 * If id == 0 : find by login
 * If login == "" : return error
 * @param  {[type]} r *Repository)  Find(ctx context.Context, id int64, login string) (user *model.User, err error [description]
 * @return {[type]}   [description]
 */
func (r *Repository) Find(ctx context.Context, id int64, login string) (user *model.User, err error) {
	query := "SELECT id, login, salt, metadata FROM %s WHERE"

	user = &model.User{}

	if id == 0 {
		query += " login = $1"
		err = r.pool.QueryRow(ctx, r.table(query), login).Scan(&user.ID, &user.Login, &user.HashedPassword, &user.Metadata)
	} else {
		query += " id = $1"
		err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&user.ID, &user.Login, &user.HashedPassword, &user.Metadata)
	}

	return
}

func (r *Repository) Update(ctx context.Context, id int64, newLogin string, newMetadata []byte) (err error) {
	const selectQuery = "SELECT login, metadata FROM %s WHERE id = $1"
	const query = "UPDATE %s SET login = $2, metadata = $3 WHERE id = $1"

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

	var (
		login string
		metadata []byte
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), id).Scan(&login, &metadata)
	if err != nil {
		return
	}

	if newLogin != "" {
		login = newLogin
	}

	if newMetadata != nil {
		metadata = newMetadata
	}

	_, err = tx.Exec(ctx, r.table(query), id, login, metadata)

	return
}

func (r *Repository) Dump(ctx context.Context) (reader io.ReadCloser, err error) {
	var (
		writer io.WriteCloser
		conn   *pgxpool.Conn
	)

	query := r.table("COPY %s TO STDOUT BINARY")

	reader, writer = io.Pipe()

	conn, err = r.pool.Acquire(ctx)
	if err != nil {
		conn.Release()
		return
	}

	// TODO: remove, see billing/invoices for example
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

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	if err != nil {
		conn.Release()
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)
	defer conn.Release()

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
