package postgres

import (
	"io"
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type TranslationsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewTranslationsRepository(pool *pgxpool.Pool, tableName string) *TranslationsRepository {
	return &TranslationsRepository{tableName: tableName, pool: pool}
}

func (r *TranslationsRepository) SaveTranslation(ctx context.Context, messageID int64, lang, title, text, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(message_id, lang, text, title, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, lang, title, text, createdAt, updatedAt)

	return
}

func (r *TranslationsRepository) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string, updatedAt string) (err error) {
	const query = "UPDATE %s SET title = $3, text = $4, updated_at = $5 WHERE message_id = $1 AND lang = $2"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, lang, title, text, updatedAt)

	return
}

func (r *TranslationsRepository) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	const query = "DELETE FROM %s WHERE message_id = $1 AND lang = $2"

	_, err = r.pool.Exec(ctx, r.table(query), messageID, lang)

	return
}

func (r *TranslationsRepository) ReadTranslation(ctx context.Context, messageID int64, lang string) (translation *model.Translation, err error) {
	const query = "SELECT title, text, created_at, updated_at FROM %s WHERE message_id = $1 AND lang = $2"

	var createdAt, updatedAt time.Time

	translation = &model.Translation{
		MessageID: messageID,
		Lang:      lang,
	}

	err = r.pool.QueryRow(ctx, r.table(query), messageID, lang).Scan(&translation.Title, &translation.Text, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	translation.CreatedAt = createdAt.Format(time.RFC3339)
	translation.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *TranslationsRepository) ListTranslations(ctx context.Context, messageID int64) (translations []*model.Translation, err error) {
	const query = "SELECT lang, title, text, created_at, updated_at FROM %s WHERE message_id = $1"

	rows, err := r.pool.Query(ctx, r.table(query), messageID)
	if err != nil {
		return nil, err
	}

	translations = make([]*model.Translation, 0)
	for rows.Next() {
		translation := &model.Translation{
			MessageID: messageID,
		}

		var createdAt, updatedAt time.Time

		err = rows.Scan(&translation.Lang, &translation.Title, &translation.Text, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		translation.CreatedAt = createdAt.Format(time.RFC3339)
		translation.UpdatedAt = updatedAt.Format(time.RFC3339)

		translations = append(translations, translation)
	}

	err = rows.Err()

	return
}

func (r *TranslationsRepository) ReadMessageTranslations(ctx context.Context, messageID int64) (previews []*model.TranslationPreview, err error) {
	const query = "SELECT lang, title, created_at, updated_at FROM %s WHERE message_id = $1"

	rows, err := r.pool.Query(ctx, r.table(query), messageID)
	if err != nil {
		return nil, err
	}

	previews = make([]*model.TranslationPreview, 0)
	for rows.Next() {
		preview := &model.TranslationPreview{
			MessageID: messageID,
		}

		var createdAt, updatedAt time.Time

		err = rows.Scan(&preview.Lang, &preview.Title, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		preview.CreatedAt = createdAt.Format(time.RFC3339)
		preview.UpdatedAt = updatedAt.Format(time.RFC3339)

		previews = append(previews, preview)
	}

	err = rows.Err()

	return
}

func (r *TranslationsRepository) DeleteMessage(ctx context.Context, messageID int64) (err error) {
	const query = "DELETE FROM %s WHERE message_id = $1"

	_, err = r.pool.Exec(ctx, r.table(query), messageID)

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