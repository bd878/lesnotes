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

func (r *Repository) SaveFile(ctx context.Context, reader io.Reader, id, userID int64, private bool, name string, mime string) (err error) {
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
			fmt.Fprintf(os.Stderr, "[SaveFile]: rollback with error: %v\n", err)
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
	defer object.Close()
	if err != nil {
		return err
	}

	var size int64
	size, err = io.Copy(object, reader)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(query), id, userID, name, private, oid, mime, size)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetMeta(ctx context.Context, id int64, fileName string) (file *model.File, err error) {
	query := "SELECT owner_id, id, name, private, oid, mime, size, created_at FROM %s WHERE"

	file = &model.File{
	}

	var createdAt time.Time

	if id != 0 {
		query += " id = $1"
		err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&file.UserID, &file.ID, &file.Name, &file.Private, &file.OID, &file.Mime, &file.Size, &createdAt)
	} else if fileName != "" {
		query += " name = $1"
		err = r.pool.QueryRow(ctx, r.table(query), fileName).Scan(&file.UserID, &file.ID, &file.Name, &file.Private, &file.OID, &file.Mime, &file.Size, &createdAt)
	} else {
		err = errors.New("id = 0 and fileName is empty")
	}
	if err != nil {
		return
	}

	file.CreateUTCNano = createdAt.UnixNano()

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
			fmt.Fprintf(os.Stderr, "[ReadFile]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	lb := tx.LargeObjects()
	object, err := lb.Open(ctx, uint32(oid), pgx.LargeObjectModeRead)
	defer object.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}

	return
}

func (r *Repository) DeleteFile(ctx context.Context, ownerID, id int64) (err error) {
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
				fmt.Fprintf(os.Stderr, "[DeleteFile]: rollback with error: %v\n", err)
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
