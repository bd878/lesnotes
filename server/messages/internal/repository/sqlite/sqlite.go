package repository

import (
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/user/pkg/model"
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

func (r *Repository) Get(ctx context.Context, userId usermodel.UserId) ([]model.Message, error) {
  rows, err := r.db.QueryContext(ctx,
    "SELECT id, message, file FROM messages WHERE user_id = ?", int(userId),
  )
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var res []model.Message
  for rows.Next() {
    var id int
    var value string
    var fileCol sql.NullString
    if err := rows.Scan(&id, &value, &fileCol); err != nil {
      return nil, err
    }

    var fileName string
    if fileCol.Valid {
      fileName = fileCol.String
    }
    res = append(res, model.Message{
      Id: id,
      UserId: -1,
      CreateTime: "null",
      Value: value,
      File: fileName,
    })
  }
  return res, nil
}
