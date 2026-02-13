package postgres

import (
	"io"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	billing "github.com/bd878/gallery/server/billing/pkg/model"
)

type PaymentsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewPaymentsRepository(pool *pgxpool.Pool, tableName string) *PaymentsRepository {
	return &PaymentsRepository{tableName, pool}
}

func (r *PaymentsRepository) SavePayment(ctx context.Context, id, userID int64, invoiceID string, currency, status string, total int64, metadata []byte) (err error) {
	const insert = "INSERT INTO %s(id, invoice_id, user_id, status, currency, total, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err = r.pool.Exec(ctx, r.table(insert), id, invoiceID, userID, status, currency, total, metadata)

	return
}

func (r *PaymentsRepository) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	const update = "UPDATE %s SET status = 'processed' WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *PaymentsRepository) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	const update = "UPDATE %s SET status = 'cancelled' WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *PaymentsRepository) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	const update = "UPDATE %s SET status = 'refunded' WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id)

	return
}

func (r *PaymentsRepository) GetPayment(ctx context.Context, id, userID int64) (payment *billing.Payment, err error) {
	const query = "SELECT invoice_id, status, currency, total, metadata FROM %s WHERE user_id = $1 AND id = $2"

	payment = &billing.Payment{
		ID:     id,
		UserID: userID,
	}

	err = r.pool.QueryRow(ctx, r.table(query), userID, id).Scan(&payment.InvoiceID, &payment.Status,
		&payment.Currency, &payment.Total, &payment.Metadata)

	return
}

func (r *PaymentsRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping payments repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *PaymentsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring payments repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r PaymentsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}