package postgres

import (
	"io"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
)

type ThreadsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewThreadsRepository(pool *pgxpool.Pool, tableName string) *ThreadsRepository {
	return &ThreadsRepository{tableName: tableName, pool: pool}
}

func (r *ThreadsRepository) Truncate(ctx context.Context) (err error) {
	logger.Debugln("truncating table")
	_, err = r.pool.Exec(ctx, r.table("TRUNCATE TABLE %s"))
	return
}

func (r *ThreadsRepository) Dump(ctx context.Context, writer io.Writer) (err error) {
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

func (r *ThreadsRepository) Restore(ctx context.Context, reader io.Reader) (err error) {
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

func (r ThreadsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}