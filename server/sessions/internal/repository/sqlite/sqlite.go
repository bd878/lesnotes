package repository

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

type Repository struct {
	pool          *sql.DB
	insertStmt    *sql.Stmt
	selectStmt    *sql.Stmt
	deleteStmt    *sql.Stmt
	deleteAllStmt *sql.Stmt
	listStmt      *sql.Stmt
}

func New(dbPath string) *Repository {
	pool, err := sql.Open("sqlite3", "file:" + dbPath)
	if err != nil {
		panic(err)
	}

	insertStmt := utils.Must(pool.Prepare(`
INSERT INTO sessions(
	user_id,
	value,
	expires_utc_nano
) VALUES (:userID, :value, :expiresUTCNano)
;`,
	))

	selectStmt := utils.Must(pool.Prepare(`
SELECT user_id, expires_utc_nano
FROM sessions
WHERE value = :token
;`,
	))

	deleteStmt := utils.Must(pool.Prepare(`
DELETE FROM sessions
WHERE value = :token
;`,
	))

	deleteAllStmt := utils.Must(pool.Prepare(`
DELETE FROM sessions
WHERE user_id = :userID
;`,
	))

	listStmt := utils.Must(pool.Prepare(`
SELECT value, expires_utc_nano
FROM sessions
WHERE user_id = :userID
;`,
	))

	return &Repository{
		pool:          pool,
		insertStmt:    insertStmt,
		selectStmt:    selectStmt,
		deleteStmt:    deleteStmt,
		listStmt:      listStmt,
		deleteAllStmt: deleteAllStmt,
	}
}

func (r *Repository) Add(ctx context.Context, userID int32, token string, expiresUTCNano int64) (err error) {
	_, err = r.insertStmt.ExecContext(ctx,
		sql.Named("userID",         userID),
		sql.Named("value",          token),
		sql.Named("expiresUTCNano", expiresUTCNano),
	)
	return err
}

func (r *Repository) Get(ctx context.Context, token string) (session *model.Session, err error) {
	var (
		expiresUTCNano int64
		userID int32
	)

	err = r.selectStmt.QueryRowContext(ctx, sql.Named("token", token)).Scan(&userID, &expiresUTCNano)
	if err != nil {
		return
	}

	session = &model.Session{
		UserID:         userID,
		Token:          token,
		ExpiresUTCNano: expiresUTCNano,
	}

	return
}

func (r *Repository) List(ctx context.Context, userID int32) (sessions []*model.Session, err error) {
	var rows *sql.Rows
	rows, err = r.listStmt.QueryContext(ctx, sql.Named("userID", userID))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions = make([]*model.Session, 0)
	for rows.Next() {
		var (
			token string
			expiresUTCNano int64
		)

		err = rows.Scan(&token, &expiresUTCNano)
		if err != nil {
			return
		}

		sessions = append(sessions, &model.Session{
			UserID:         userID,
			Token:          token,
			ExpiresUTCNano: expiresUTCNano,
		})
	}

	return
}

func (r *Repository) Delete(ctx context.Context, token string) (err error) {
	_, err = r.deleteStmt.ExecContext(ctx, sql.Named("token", token))
	return
}

func (r *Repository) DeleteAll(ctx context.Context, userID int32) (err error) {
	_, err = r.deleteAllStmt.ExecContext(ctx, sql.Named("userID", userID))
	return
}