package postgres

import (
	"fmt"
	"io"
	"os"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type UsersRepository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool, tableName string) *UsersRepository {
	return &UsersRepository{
		tableName: tableName,
		pool:      pool,
	}
}

func (r *UsersRepository) Save(ctx context.Context, id int64, login, salt string, metadata []byte) (err error) {
	const query = "INSERT INTO %s(id, login, salt, metadata) VALUES ($1, $2, $3, $4)"

	_, err = r.pool.Exec(ctx, r.table(query), id, login, salt, metadata)

	return
}

func (r *UsersRepository) Delete(ctx context.Context, id int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), id)

	return
}

/**
 * Find by id or login
 * If id == 0 : find by login
 * If login == "" : return error
 * @param  {[type]} r *UsersRepository)  Find(ctx context.Context, id int64, login string) (user *model.User, err error [description]
 * @return {[type]}   [description]
 */
func (r *UsersRepository) Find(ctx context.Context, id int64, login string) (user *model.User, err error) {
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

func (r *UsersRepository) Update(ctx context.Context, id int64, newLogin string, newMetadata []byte) (err error) {
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

func (r *UsersRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping users repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *UsersRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring users repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r UsersRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
