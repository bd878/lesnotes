package postgres

import (
	"io"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/search/pkg/model"
)

type TranslationsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewTranslationsRepository(pool *pgxpool.Pool, tableName string) *TranslationsRepository {
	return &TranslationsRepository{tableName: tableName, pool: pool}
}

func (r *TranslationsRepository) SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error) {
	const query = "INSERT INTO %s(message_id, user_id, lang, title, text) VALUES ($1, $2, $3, $4, $5)"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, userID, lang, title, text)

	return
}

func (r *TranslationsRepository) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	const query = "DELETE FROM %s WHERE message_id = $1 AND lang = $2"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, lang)

	return
}

func (r *TranslationsRepository) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	const query = "UPDATE %s SET title = $3, text = $4 WHERE message_id = $1 AND lang = $2"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, lang, title, text)

	return
}

func (r *TranslationsRepository) SearchTranslations(ctx context.Context, userID int64, substr string) (list []*model.Translation, err error) {
	const query = "SELECT message_id, lang, text, title FROM %s WHERE user_id = $1 AND text || ' ' || title ILIKE $2"

	rows, err := r.pool.Query(ctx, r.table(query), userID, "'%" + substr + "%'")
	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*model.Translation, 0)
	for rows.Next() {
		translation := &model.Translation{}

		err = rows.Scan(&translation.MessageID, &translation.Lang, &translation.Text, &translation.Title)
		if err != nil {
			return
		}

		list = append(list, translation)
	}

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