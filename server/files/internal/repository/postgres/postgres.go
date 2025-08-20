package repository

import (
	"fmt"
	"os"
	"io"
	"time"
	"errors"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Repository struct {
	tableName  string
	pool      *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		tableName: "files.files",
		pool:      pool,
	}
}

func (r *Repository) SaveFile(ctx context.Context, ownerID int32, reader io.Reader, file *model.File) (err error) {
	const query = "INSERT INTO %s(id, owner_id, name, private, oid, mime, size) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	const createdAtQuery = "SELECT created_at FROM %s WHERE id = $1"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var oid uint32
	lb := tx.LargeObjects()
	oid, err = lb.Create(ctx, 0)
	if err != nil {
		logger.Errorw("failed to create large object", "error", err)
		return err
	}

	object, err := lb.Open(ctx, oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return err
	}
	defer object.Close()

	var size int64
	size, err = io.Copy(object, reader)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(query), file.ID, file.UserID, file.Name, file.Private, oid, file.Mime, size)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetMeta(ctx context.Context, ownerID int32, id int32) (file *model.File, err error) {
	const query = "SELECT name, private, oid, mime, size, created_at FROM %s WHERE owner_id = $1 AND id = $2"

	logger.Infow("get meta", "owner_id", ownerID, "id", id)

	file = &model.File{
		UserID: ownerID,
		ID:     id,
	}

	var created time.Time

	err = r.pool.QueryRow(ctx, r.table(query), ownerID, id).Scan(&file.Name, &file.Private, &file.OID, &file.Mime, &file.Size, &created)
	if err != nil {
		return
	}

	file.CreateUTCNano = created.UnixNano()

	return
}

func (r *Repository) ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error) {
	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	lb := tx.LargeObjects()
	object, err := lb.Open(ctx, uint32(oid), pgx.LargeObjectModeRead)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}

	object.Close()

	return
}

func (r *Repository) DeleteFile(ctx context.Context, ownerID, id int32) (err error) {
	const query = "SELECT oid FROM %s WHERE owner_id = $1 AND id = $2"
	const deleteQuery = "DELETE FROM %s WHERE owner_id = $1 AND id = $2"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			if errors.Is(err, pgx.ErrNoRows) {
				return
			} else {
				fmt.Fprintf(os.Stderr, "rollback with error: %v\n", err)
				err = tx.Rollback(ctx)
			}
		default:
			err = tx.Commit(ctx)
		}
	}()

	var oid int
	err = tx.QueryRow(ctx, r.table(query), ownerID, id).Scan(&oid)
	if err != nil {
		return
	}

	lb := tx.LargeObjects()
	err = lb.Unlink(ctx, uint32(oid))
	if err != nil {
		return
	}

	result, err := tx.Exec(ctx, r.table(deleteQuery), ownerID, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("no rows owner_id %d id %d\n", ownerID, id)
	}

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
