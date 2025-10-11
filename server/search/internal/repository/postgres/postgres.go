package postgres

import (
	"context"

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
	return nil
}

func (r *Repository) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	return nil
}

