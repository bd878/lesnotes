package repository

import (
	"fmt"
	"os"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/users/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		tableName: "users.users",
		pool:      pool,
	}
}

func (r *Repository) Save(ctx context.Context, id int64, login, salt, theme string) (err error) {
	const query = "INSERT INTO %s(id, login, salt, theme) VALUES ($1, $2, $3, $4)"

	_, err = r.pool.Exec(ctx, r.table(query), id, login, salt, theme)

	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1"

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

	_, err = tx.Exec(ctx, r.table(query), id)

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
	query := "SELECT id, login, salt, theme FROM %s WHERE"

	user = &model.User{}

	if id == 0 {
		query += " login = $1"
		err = r.pool.QueryRow(ctx, r.table(query), login).Scan(&user.ID, &user.Login, &user.HashedPassword, &user.Theme)
	} else {
		query += " id = $1"
		err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&user.ID, &user.Login, &user.HashedPassword, &user.Theme)
	}

	return
}

func (r *Repository) Update(ctx context.Context, id int64, newLogin, newTheme string) (err error) {
	const selectQuery = "SELECT login, theme FROM %s WHERE id = $1"
	const query = "UPDATE %s SET login = $2, theme = $3 WHERE id = $1"

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
		login, theme string
	)

	err = tx.QueryRow(ctx, r.table(selectQuery), id).Scan(&login, &theme)
	if err != nil {
		return
	}

	if newLogin != "" {
		login = newLogin
	}

	if newTheme != "" {
		theme = newTheme
	}

	_, err = tx.Exec(ctx, r.table(query), id, login, theme)

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}