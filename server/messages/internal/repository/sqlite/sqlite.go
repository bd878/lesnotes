package repository

import (
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository struct {
  db *sql.DB
}

func New(dbpath string) (*Repository, error) {
  db, err := sql.Open("sqlite3", "file:" + dbpath)
  if err != nil {
    return nil, err
  }
  return &Repository{db}, nil
}

func (r *Repository) Put(ctx context.Context, msg *model.Message) error {
  _, err := r.db.ExecContext(ctx,
    "INSERT INTO messages(user_id, createtime, message, file) VALUES (?,?,?,?)",
    msg.UserId, msg.CreateTime, msg.Value, msg.File,
  )
  return err
}

func (r *Repository) GetAll(ctx context.Context) ([]model.Message, error) {
  rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, message, file FROM messages")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var res []model.Message
  for rows.Next() {
    var id, userId int
    var value string
    var fileCol sql.NullString
    if err := rows.Scan(&id, &userId, &value, &fileCol); err != nil {
      return nil, err
    }

    var fileName string
    if fileCol.Valid {
      fileName = fileCol.String
    }
    res = append(res, model.Message{
      Id: id,
      UserId: userId,
      Value: value,
      File: fileName,
    })
  }
  return res, nil
}
