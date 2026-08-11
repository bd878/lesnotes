package postgres

import (
	"fmt"
	"errors"
	"log/slog"
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/tm"
)

type OutboxStore struct {
	db DB
	tableName string
}

var _ tm.OutboxStore = (*OutboxStore)(nil)

func NewOutboxStore(db DB, tableName string) OutboxStore {
	return OutboxStore{
		tableName: tableName,
		db: db,
	}
}

func (s OutboxStore) Save(ctx context.Context, msg am.Message) error {
	const query = "INSERT INTO %s(id, name, subject, data) VALUES ($1, $2, $3, $4)"

	_, err := s.db.Exec(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				slog.Info("unique violation, outbox message exist",
					slog.String("message_name", msg.MessageName()),
					slog.String("subject", msg.Subject()),
				)
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
	}

	return err
}

func (s OutboxStore) FindUnpublished(ctx context.Context, limit int) ([]am.Message, error) {
	return nil, nil
}

func (s OutboxStore) MarkPublished(ctx context.Context, ids ...string) error {
	return nil
}

func (s OutboxStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
