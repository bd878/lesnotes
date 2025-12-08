package postgres

import (
	"io"
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
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

func (m *PaymentsRepository) SavePayment(ctx context.Context, id, userID int64, invoiceID string, currency, status string, total int64, metadata []byte) (err error) {
	return
}

func (m *PaymentsRepository) ProceedPayment(ctx context.Context, id, userID int64) (err error) {
	return
}

func (m *PaymentsRepository) CancelPayment(ctx context.Context, id, userID int64) (err error) {
	return
}

func (m *PaymentsRepository) RefundPayment(ctx context.Context, id, userID int64) (err error) {
	return
}

func (m *PaymentsRepository) GetPayment(ctx context.Context, id, userID int64) (payment *billing.Payment, err error) {
	return
}

func (r *PaymentsRepository) Truncate(ctx context.Context) (err error) {
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

func (r *PaymentsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
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

func (r PaymentsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}