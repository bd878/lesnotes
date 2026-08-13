package postgres

import (
	"fmt"
	"errors"
	"log/slog"
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bd878/gallery/server/internal/tm"
	"github.com/bd878/gallery/server/internal/am"
)

type InboxStore struct {
	db DB
	tableName string
}

var _ tm.InboxStore = (*InboxStore)(nil)

func NewInboxStore(db DB, tableName string) InboxStore {
	return InboxStore{
		tableName: tableName,
		db: db,
	}
}

func (s InboxStore) Save(ctx context.Context, msg am.RawMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data, received_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)"

	_, err := s.db.Exec(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				slog.Info("unique violation, inbox message exist",
					slog.String("message_name", msg.MessageName()),
					slog.String("subject", msg.Subject()),
				)
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
	}

	return err
}

func (s InboxStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}