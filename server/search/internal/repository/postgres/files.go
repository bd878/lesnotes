package postgres

import (
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FilesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewFilesRepository(pool *pgxpool.Pool, tableName string) *FilesRepository {
	return &FilesRepository{tableName: tableName, pool: pool}
}

func (r *FilesRepository) SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, owner_id, name, description, mime, size, private, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $0)"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID, name, description,
		mime, size, private, createdAt, updatedAt)

	return
}

func (r *FilesRepository) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1 AND owner_id = $2"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID)

	return
}

func (r *FilesRepository) PublishFile(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = false, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)

	return
}

func (r *FilesRepository) PrivateFile(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)

	return
}

func (r FilesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}