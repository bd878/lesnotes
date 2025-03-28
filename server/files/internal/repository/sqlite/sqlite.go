package repository

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Repository struct {
	pool          *sql.DB
	insertStmt    *sql.Stmt
	selectStmt    *sql.Stmt
}

func New(dbPath string) *Repository {
	pool, err := sql.Open("sqlite3", "file:" + dbPath)
	if err != nil {
		panic(err)
	}

	insertStmt := utils.Must(pool.Prepare(`
INSERT INTO files(
	id,
	user_id,
	name,
	create_utc_nano,
	private
) VALUES (:id, :userId, :name, :createUtcNano, :private)
;`,
	))

	selectStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, name, create_utc_nano, private
FROM files
WHERE id = :id AND user_id = :userId
;`,
	))

	return &Repository{
		pool: pool,
		insertStmt: insertStmt,
		selectStmt: selectStmt,
	}
}

func (r *Repository) SaveFile(ctx context.Context, log *logger.Logger, file *model.File) error {
	var privateCol sql.NullInt32
	if file.Private {
		privateCol.Int32 = 1
		privateCol.Valid = true
	}

	_, err := r.insertStmt.ExecContext(ctx,
		sql.Named("id", file.ID),
		sql.Named("userId", file.UserID),
		sql.Named("name", file.Name),
		sql.Named("createUtcNano", file.CreateUTCNano),
		sql.Named("private", privateCol),
	)
	if err != nil {
		log.Errorw("failed to save file", "id", file.ID, "user_id", file.UserID)
		return err
	}
	return nil
}

func (r *Repository) ReadFile(ctx context.Context, log *logger.Logger, params *model.ReadFileParams) (
	*model.File, error,
) {
	var (
		id, userId int32
		name string
		createUTCNano int64
		privateCol sql.NullInt32
	)

	err := r.selectStmt.QueryRowContext(ctx, sql.Named("id", params.ID),
		sql.Named("userId", params.UserID)).Scan(&id, &userId, &name, &createUTCNano, &privateCol)

	msg := &model.File{
		ID:                id,
		UserID:            userId,
		Name:              name,
		CreateUTCNano:     createUTCNano,
		Private:           true,
	}

	if privateCol.Valid {
		if privateCol.Int32 == 0 {
			msg.Private = false
		}
	}

	switch {
	case err == sql.ErrNoRows:
		log.Errorw("failed to read file, no rows found", "id", params.ID, "user_id", params.UserID)
		return nil, err

	case err != nil:
		log.Errorw("failed to read file, unknown error", "id", params.ID, "user_id", params.UserID)
		return nil, err

	default:
		return msg, nil
	}
}