package postgres

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Dumper struct {
	pool              *pgxpool.Pool
	invoicesTableName string
	paymentsTableName string
	ctx               context.Context
	cancel            context.CancelCauseFunc
	ch                chan *api.BillingSnapshot
	wg                sync.WaitGroup
}

func NewDumper(pool *pgxpool.Pool, invoicesTableName, paymentsTableName string) *Dumper {
	return &Dumper{
		pool:                 pool,
		invoicesTableName:    invoicesTableName,
		paymentsTableName:    paymentsTableName,
	}
}

func (r *Dumper) Open(ctx context.Context) (ch chan *api.BillingSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.BillingSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.dumpInvoices()
	r.wg.Add(1)
	go r.dumpPayments()

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()

	return
}

func (r *Dumper) dumpInvoices() {
	query := "SELECT id, user_id, status, total, cart, metadata, created_at, updated_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("invoices dump finished")

	rows, err := r.pool.Query(r.ctx, r.invoicesTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		invoice := &api.InvoiceSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&invoice.Id, &invoice.UserId, &invoice.Status, &invoice.Total,
			&invoice.Cart, &invoice.Metadata, &createdAt, &updatedAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		invoice.CreatedAt = createdAt.Format(time.RFC3339)
		invoice.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.BillingSnapshot{
			Item: &api.BillingSnapshot_Invoice{
				Invoice: invoice,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) dumpPayments() {
	query := "SELECT id, invoice_id, user_id, status, currency, total, metadata, created_at, updated_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("payments dump finished")

	rows, err := r.pool.Query(r.ctx, r.paymentsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		payment := &api.PaymentSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&payment.Id, &payment.InvoiceId, &payment.UserId, &payment.Status,
			&payment.Currency, &payment.Total, &payment.Metadata, &createdAt, &updatedAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		payment.CreatedAt = createdAt.Format(time.RFC3339)
		payment.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.BillingSnapshot{
			Item: &api.BillingSnapshot_Payment{
				Payment: payment,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) Close() (err error) {
	logger.Debugln("close dumper")
	r.cancel(nil)
	r.wg.Wait()
	return nil
}

func (r *Dumper) Restore(ctx context.Context, snapshot *api.BillingSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.BillingSnapshot_Invoice:

		query := "INSERT INTO %s(id, user_id, status, total, cart, metadata, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

		_, err = r.pool.Exec(ctx, r.invoicesTable(query), v.Invoice.Id, v.Invoice.UserId, v.Invoice.Status,
			v.Invoice.Total, v.Invoice.Cart, v.Invoice.Metadata, v.Invoice.CreatedAt, v.Invoice.UpdatedAt)

		return

	case *api.BillingSnapshot_Payment:

		query := "INSERT INTO %s(id, invoice_id, user_id, status, currency, total, metadata, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"

		_, err = r.pool.Exec(ctx, r.paymentsTable(query), v.Payment.Id, v.Payment.InvoiceId, v.Payment.UserId, v.Payment.Status,
			v.Payment.Currency, v.Payment.Total, v.Payment.Metadata, v.Payment.CreatedAt, v.Payment.UpdatedAt)

		return

	default:
		return errors.New("unknown snapshot item")
	}
}

func (r Dumper) invoicesTable(query string) string {
	return fmt.Sprintf(query, r.invoicesTableName)
}

func (r Dumper) paymentsTable(query string) string {
	return fmt.Sprintf(query, r.paymentsTableName)
}
