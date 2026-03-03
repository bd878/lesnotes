package postgres

import (
	"fmt"
	"time"
	"io"
	"context"

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

func (r *UsersRepository) Save(ctx context.Context, id int64, login, salt string, metadata []byte, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, login, salt, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.table(query), id, login, salt, metadata, createdAt, updatedAt)

	return
}

func (r *UsersRepository) Delete(ctx context.Context, id int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), id)

	return
}

func (r *UsersRepository) FindByID(ctx context.Context, id int64) (user *model.User, err error) {
	query := "SELECT login, salt, metadata, created_at, updated_at FROM %s WHERE id = $1"

	user = &model.User{
		ID:     id,
	}

	var createdAt, updatedAt *time.Time

	err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&user.Login, &user.HashedPassword, &user.Metadata, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *UsersRepository) FindByLogin(ctx context.Context, login string) (user *model.User, err error) {
	query := "SELECT id, salt, metadata, created_at, updated_at FROM %s WHERE login = $1"

	user = &model.User{
		Login:   login,
	}

	var createdAt, updatedAt *time.Time

	err = r.pool.QueryRow(ctx, r.table(query), login).Scan(&user.ID, &user.HashedPassword, &user.Metadata, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *UsersRepository) Update(ctx context.Context, id int64, login *string, metadata []byte, updatedAt string) (err error) {
	const query = "UPDATE %s SET login = $2, metadata = $3, updated_at = $4 WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), id, login, metadata, updatedAt)

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
