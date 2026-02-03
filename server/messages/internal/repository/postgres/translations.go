package postgres

import (
	"io"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type TranslationsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewTranslationsRepository(tableName string, pool *pgxpool.Pool) *TranslationsRepository {
	return &TranslationsRepository{tableName: tableName, pool: pool}
}

func (r *TranslationsRepository) SaveTranslation(ctx context.Context, messageID int64, lang, text, title string) (err error) {
	return
}

func (r *TranslationsRepository) UpdateTranslation(ctx context.Context, messageID int64, lang string, text, title *string) (err error) {
	return
}

func (r *TranslationsRepository) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	return
}

func (r *TranslationsRepository) ReadTranslation(ctx context.Context, messageID int64, lang string) (translation *model.Translation, err error) {
	return
}

func (r *TranslationsRepository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.table("TRUNCATE TABLE %s"))
	return
}

func (r *TranslationsRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping translations repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *TranslationsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring translations repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r TranslationsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}