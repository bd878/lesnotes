package postgres

import (
	"io"
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type InvoicesRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewInvoicesRepository(pool *pgxpool.Pool, tableName string) *InvoicesRepository {
	return &InvoicesRepository{tableName: tableName, pool: pool}
}

func (r *InvoicesRepository) SaveInvoice(ctx context.Context, id string, userID int64, currency, status string, total int64,
	metadata []byte, cart []byte, createdAt, updatedAt string) (err error) {
	const insert = "INSERT INTO %s(id, user_id, currency, status, total, metadata, cart, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err = r.pool.Exec(ctx, r.table(insert), id, userID, currency, status, total, metadata, cart, createdAt, updatedAt)

	return
}

func (r *InvoicesRepository) PayInvoice(ctx context.Context, id string, userID int64, updatedAt string) (err error) {
	const update = "UPDATE %s SET status = 'paid', updated_at = $3 WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id, updatedAt)

	return
}

func (r *InvoicesRepository) CancelInvoice(ctx context.Context, id string, userID int64, updatedAt string) (err error) {
	const update = "UPDATE %s SET status = 'cancel', updated_at = $3 WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id, updatedAt)

	return
}

func (r *InvoicesRepository) GetInvoice(ctx context.Context, id string, userID int64) (invoice *api.Invoice, err error) {
	const query = "SELECT status, currency, total, metadata, cart, created_at, updated_at FROM %s WHERE user_id = $1 AND id = $2"

	invoice = &api.Invoice{
		Id:      id,
		UserId:  userID,
	}

	var createdAt, updatedAt *time.Time

	var cart []byte

	err = r.pool.QueryRow(ctx, r.table(query), userID, id).Scan(&invoice.Status, &invoice.Currency,
		&invoice.Total, &invoice.Metadata, &cart, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	var cc api.Cart
	err = proto.Unmarshal(cart, &cc)
	if err != nil {
		return
	}

	invoice.Cart = &cc
	invoice.CreatedAt = createdAt.Format(time.RFC3339)
	invoice.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *InvoicesRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("dumping invoices repo")

	conn, err = r.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return
	}

	// will block, not concurrent safe
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, r.table("COPY %s TO STDOUT BINARY"))

	return
}

func (r *InvoicesRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
	var conn *pgxpool.Conn

	logger.Debugln("restoring invoices repo")

	query := r.table("COPY %s FROM STDIN BINARY")

	conn, err = r.pool.Acquire(ctx) 
	defer conn.Release()
	if err != nil {
		return
	}

	_, err = conn.Conn().PgConn().CopyFrom(ctx, reader, query)

	return
}

func (r InvoicesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}