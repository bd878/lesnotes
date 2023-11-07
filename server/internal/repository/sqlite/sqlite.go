package repository

import (
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/pkg/model"
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

func (r *Repository) Put(ctx context.Context, message string) error {
  _, err := r.db.ExecContext(ctx, "INSERT INTO messages(message) VALUES (?)", message)
  return err
}

func (r *Repository) GetAll(ctx context.Context) ([]model.Message, error) {
  rows, err := r.db.QueryContext(ctx, "SELECT message FROM messages")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var res []model.Message
  for rows.Next() {
    var message string
    if err := rows.Scan(&message); err != nil {
      return nil, err
    }
    res = append(res, model.Message{
      Value: message,
    })
  }
  return res, nil
}