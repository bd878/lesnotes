package postgres

import (
	"io"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	billing "github.com/bd878/gallery/server/billing/pkg/model"
)

type InvoicesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewInvoicesRepository(pool *pgxpool.Pool, tableName string) *InvoicesRepository {
	return &InvoicesRepository{tableName: tableName, pool: pool}
}

func (r *InvoicesRepository) SaveInvoice(ctx context.Context, id string, userID int64, currency, status string, total int64, metadata []byte) (err error) {
	const insert = "INSERT INTO %s(id, user_id, currency, status, total, metadata) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.table(insert), id, userID, currency, status, total, metadata)

	return
}

func (r *InvoicesRepository) PayInvoice(ctx context.Context, id string, userID int64) (err error) {
	const update = "UPDATE %s SET status = 'paid' WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *InvoicesRepository) CancelInvoice(ctx context.Context, id string, userID int64) (err error) {
	const update = "UPDATE %s SET status = 'cancel' WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *InvoicesRepository) GetInvoice(ctx context.Context, id string, userID int64) (invoice *billing.Invoice, err error) {
	const query = "SELECT status, currency, total, metadata FROM %s WHERE user_id = $1 AND id = $2"

	invoice = &billing.Invoice{
		ID:      id,
		UserID:  userID,
	}

	err = r.pool.QueryRow(ctx, r.table(query), userID, id).Scan(&invoice.Status, &invoice.Currency,
		&invoice.Total, &invoice.Metadata)

	return
}

func (r *InvoicesRepository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.table("TRUNCATE TABLE %s"))
	return
}

func (r *InvoicesRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	conn, err = r.pool.Acquire(ctx)
	if err != nil {
		conn.Release()
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))
	defer conn.Release()
	if err != nil {
		return
	}

	return
}

func (r *InvoicesRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	if err != nil {
		conn.Release()
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)
	defer conn.Release()

	return
}

func (r InvoicesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}