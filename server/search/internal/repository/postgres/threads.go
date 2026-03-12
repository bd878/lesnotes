package postgres

import (
	"os"
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/search/pkg/model"
)

type ThreadsRepository struct {
	tableName        string
	pool             *pgxpool.Pool
}

func NewThreadsRepository(pool *pgxpool.Pool, tableName string) *ThreadsRepository {
	return &ThreadsRepository{tableName: tableName, pool: pool}
}

func (r *ThreadsRepository) SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, user_id, parent_id, name, description, private, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err = r.pool.Exec(ctx, r.table(query), id, userID, parentID, name, description, private, createdAt, updatedAt)

	return
}

func (r *ThreadsRepository) SearchThreads(ctx context.Context, parentID, userID int64) (list []*model.Thread, err error) {
	const query = "SELECT id, user_id, description, private, created_at, updated_at FROM %s WHERE user_id = $1 AND parent_id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[SearchThreads]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows

	rows, err = r.pool.Query(ctx, r.table(query), userID, parentID)
	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*model.Thread, 0)
	for rows.Next() {
		var createdAt, updatedAt *time.Time

		thread := &model.Thread{
			UserID: userID,
		}

		err = rows.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.Private, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		thread.CreatedAt = createdAt.Format(time.RFC3339)
		thread.UpdatedAt = updatedAt.Format(time.RFC3339)

		list = append(list, thread)
	}

	return
}

func (r *ThreadsRepository) UpdateThread(ctx context.Context, id, userID int64, name, description *string, updatedAt string) (err error) {
	const update = "UPDATE %s SET name = $3, description = $4, updated_at = $5 WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id, name, description, updatedAt)

	return
}

func (r *ThreadsRepository) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	const query = "DELETE FROM %s WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(query), userID, id)

	return
}

func (r *ThreadsRepository) ChangeThreadParent(ctx context.Context, id, userID, newParentID int64) (err error) {
	const query = "SELECT parent_id FROM %s WHERE user_id = $1 AND id = $2"
	const update = "UPDATE %s SET parent_id = $3 WHERE user_id = $1 AND id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[ChangeThreadParent]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var parentID int64

	if newParentID != 0 {
		parentID = newParentID
	}

	err = tx.QueryRow(ctx, r.table(query), userID, id).Scan(&parentID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(update), userID, id, parentID)

	return
}

func (r *ThreadsRepository) PublishThread(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	const update = "UPDATE %s SET private = false, updated_at = $3 WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id, updatedAt)

	return
}

func (r *ThreadsRepository) PrivateThread(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	const update = "UPDATE %s SET private = true, updated_at = $3 WHERE user_id = $1 AND id = $2"

	_, err = r.pool.Exec(ctx, r.table(update), userID, id, updatedAt)

	return
}

func (r ThreadsRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}