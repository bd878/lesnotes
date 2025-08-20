package repository

import (
	"fmt"
	"context"
	"errors"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/repository"
)

// TODO: remove, stale
type Repository struct {
	tableName   string
	pool        *sql.DB
}

func New(tableName, dbPath string) *Repository {
	pool, err := sql.Open("sqlite3", "file:" + dbPath)
	if err != nil {
		panic(err)
	}

	return &Repository{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *Repository) Add(ctx context.Context, id int32, name, password string) (err error) {
	const query = "INSERT INTO %s(id, name, password VALUES(?,?,?)"

	_, err = r.pool.ExecContext(ctx, r.table(query), id, name, password)

	return
}

func (r *Repository) Find(ctx context.Context, params *model.FindUserParams) (user *model.User, err error) {
	const query = "SELECT id, name, password FROM %s WHERE name = :name"

	var (
		id int32
		name, password string
	)

	err = r.pool.QueryRowContext(ctx, r.table(query), sql.Named("name", params.Name)).Scan(&id, &name, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return
	}

	user = &model.User{
		ID:       id,
		Name:     name,
		Password: password,
	}

	return
}

func (r *Repository) Get(ctx context.Context, id int32) (user *model.User, err error) {
	const query = "SELECT name, password FROM %s WHERE id = :id"

	var name, password string

	err = r.pool.QueryRowContext(ctx, r.table(query), sql.Named("id", id)).Scan(&name, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return
	}

	user = &model.User{
		ID:        id,
		Name:      name,
		Password:  password,
	}

	return
}

func (r *Repository) Delete(ctx context.Context, id int32) (err error) {
	const query = "DELETE FROM %s WHERE id = :id"

	var tx *sql.Tx

	tx, err = r.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback()
			panic(p)
		case err != nil:
			rErr := tx.Rollback()
			if rErr != nil {
				logger.Default().Errorw("failed to rollback delete stmt", "error", rErr)
			}
		default:
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, r.table(query), sql.Named("id", id))

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}