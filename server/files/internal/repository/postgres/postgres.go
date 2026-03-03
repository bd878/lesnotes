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

	"github.com/bd878/gallery/server/internal/logger"
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

func (r *Repository) SaveFile(ctx context.Context, reader io.Reader, id, userID int64, private bool, name, description, mime, createdAt, updatedAt string) (size int64, err error) {
	const query = "INSERT INTO %s(id, owner_id, name, description, private, oid, mime, size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	const createdAtQuery = "SELECT created_at FROM %s WHERE id = $1"

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
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
		return
	}

	object, err := lb.Open(ctx, oid, pgx.LargeObjectModeWrite)
	defer object.Close()
	if err != nil {
		return 0, err
	}

	size, err = io.Copy(object, reader)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, r.table(query), id, userID, name, description, private, oid, mime, size, createdAt, updatedAt)

	return
}

func (r *Repository) GetMetaByID(ctx context.Context, id int64) (file *model.File, err error) {
	query := "SELECT owner_id, name, description, private, oid, mime, size, created_at, updated_at FROM %s WHERE id = $1"

	file = &model.File{
		ID:   id,
	}

	var createdAt, updatedAt *time.Time
	err = r.pool.QueryRow(ctx, r.table(query), id).Scan(&file.UserID, &file.Name, &file.Description,
		&file.Private, &file.OID, &file.Mime, &file.Size, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	file.CreatedAt = createdAt.Format(time.RFC3339)
	file.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *Repository) GetMetaByName(ctx context.Context, fileName string) (file *model.File, err error) {
	query := "SELECT owner_id, id, description, private, oid, mime, size, created_at, updated_at FROM %s WHERE name = $1"

	file = &model.File{
		Name:  fileName,
	}

	var createdAt, updatedAt *time.Time
	err = r.pool.QueryRow(ctx, r.table(query), fileName).Scan(&file.UserID, &file.ID,
		&file.Description, &file.Private, &file.OID, &file.Mime, &file.Size, &createdAt, &updatedAt)
	if err != nil {
		return
	}

	file.CreatedAt = createdAt.Format(time.RFC3339)
	file.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}


func (r *Repository) ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*model.File, isLastPage bool, err error) {
	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, false, err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			fmt.Fprintf(os.Stderr, "[ListFiles]: rollback with error: %v\n", err)
			err = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	var rows pgx.Rows

	// TODO: handle ascending field
	if private {
		// list all : private and public
		rows, err = tx.Query(ctx, r.table("SELECT id, name, description, private, mime, size, created_at, updated_at FROM %s WHERE owner_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"), userID, limit, offset)
	} else {
		// list public only (not private)
		rows, err = tx.Query(ctx, r.table("SELECT id, name, description, private, mime, size, created_at, updated_at FROM %s WHERE owner_id = $1 AND private = false ORDER BY created_at DESC LIMIT $2 OFFSET $3"), userID, limit, offset)
	}

	defer rows.Close()
	if err != nil {
		return
	}

	list = make([]*model.File, 0)
	for rows.Next() {
		var createdAt, updatedAt time.Time

		file := &model.File{
			UserID:   userID,
		}

		err = rows.Scan(&file.ID, &file.Name, &file.Description, &file.Private, &file.Mime, &file.Size, &createdAt, &updatedAt)
		if err != nil {
			return
		}

		file.CreatedAt = createdAt.Format(time.RFC3339)
		file.UpdatedAt = updatedAt.Format(time.RFC3339)

		list = append(list, file)
	}

	if err = rows.Err(); err != nil {
		return
	}

	if int32(len(list)) < limit {
		isLastPage = true
	} else {
		var count int32
		err = tx.QueryRow(ctx, r.table("SELECT COUNT(*) FROM %s WHERE owner_id = $1"), userID).Scan(&count)
		if err != nil {
			return
		}

		if count <= offset + limit {
			isLastPage = true
		}
	}

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

func (r *Repository) PrivateFile(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = true, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)

	return
}

func (r *Repository) PublishFile(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	_, err = r.pool.Exec(ctx, r.table("UPDATE %s SET private = false, updated_at = $3 WHERE owner_id = $1 AND id = $2"), userID, id, updatedAt)

	return
}

func (r Repository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
