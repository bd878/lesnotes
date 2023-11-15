package repository

import (
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/user/pkg/model"
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

func (r *Repository) Add(ctx context.Context, usr *model.User) error {
  _, err := r.db.ExecContext(ctx, "INSERT INTO users(name,password,token)" +
    "VALUES(?,?,?)", usr.Name, usr.Password, usr.Token)
  return err
}