package postgres

import (
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
)

type FilesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewFilesRepository(pool *pgxpool.Pool, tableName string) *FilesRepository {
	return &FilesRepository{tableName: tableName, pool: pool}
}

func (r *FilesRepository) SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, owner_id, name, description, mime, size, private, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID, name, description,
		mime, size, private, createdAt, updatedAt)

	return
}

func (r *FilesRepository) DeleteFiles(ctx context.Context, ids []int64, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1 AND owner_id = $2"

	for _, id := range ids {
		_, err = r.pool.Exec(ctx, r.table(query), id, userID)
		if err != nil {
			logger.Errorln(err)
		}
	}

	return
}

func (r *FilesRepository) PublishFiles(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	for _, id := range ids {
		_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = false, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)
		if err != nil {
			logger.Errorln(err)
		}
	}

	return
}

func (r *FilesRepository) PrivateFiles(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	for _, id := range ids {
		_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)
		if err != nil {
			logger.Errorln(err)
		}
	}

	return
}

func (r FilesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}