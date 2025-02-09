package repository

import (
  "errors"
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/files/pkg/model"
)

type Repository struct {
  db          *sql.DB
}

func New(dbPath string) *Repository {
  db, err := sql.Open("sqlite3", "file:" + dbPath)
  if err != nil {
    panic(err)
  }
  return &Repository{db}
}

func (r *Repository) SaveFile(ctx context.Context, log *logger.Logger, file *model.File) error {
  _, err := r.db.ExecContext(ctx, "INSERT INTO files(id,user_id,name,create_utc_nano)" +
    "VALUES(?,?,?,?)", file.ID, file.UserID, file.Name, file.CreateUTCNano)
  if err != nil {
    log.Error("query error: %v\n", err)
  }
  return err
}

func (r *Repository) ReadFile(ctx context.Context, log *logger.Logger, params *model.ReadFileParams) (
  *model.File, error,
) {
  var id int32
  var name string
  var createUTCNano int64

  err := r.db.QueryRowContext(ctx, "SELECT id, name, create_utc_nano FROM files WHERE " +
    "id = ? AND user_id = ?", params.ID, params.UserID).Scan(&id, &name, &createUTCNano)

  msg := &model.File{
    ID:                id,
    UserID:            params.UserID,
    Name:              name,
    CreateUTCNano:     createUTCNano,
  }

  switch {
  case err == sql.ErrNoRows:
    log.Error("no rows for name %v\n", name)
    return nil, errors.New("no user")

  case err != nil:
    log.Error("query error: %v\n", err)
    return nil, err

  default:
    return msg, nil
  }
}