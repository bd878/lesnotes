package repository

import (
  "context"
  "log"
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

func (r *Repository) Add(ctx context.Context, user *model.User) error {
  _, err := r.db.ExecContext(ctx, "INSERT INTO users(name,password,token,expires)" +
    "VALUES(?,?,?,?)", user.Name, user.Password, user.Token, user.Expires)
  if err != nil {
    log.Printf("query error: %v\n", err)
  }
  return err
}

func (r *Repository) Has(ctx context.Context, user *model.User) (bool, error) {
  if user.Password == "" {
    return r.hasUser(ctx, user.Name)
  }

  return r.hasUserAndPassword(ctx, user)
}

func (r *Repository) Refresh(ctx context.Context, user *model.User) error {
  _, err := r.db.ExecContext(ctx, "UPDATE users SET token = ?, expires = ? WHERE name = ?",
    user.Token, user.Expires, user.Name)
  return err
}

func (r *Repository) hasUserAndPassword(ctx context.Context, user *model.User) (bool, error) {
  var count int
  err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users WHERE " +
    "name = ? AND password = ?", user.Name, user.Password).Scan(&count)
  switch {
  case err != nil:
    log.Printf("query error: %v\n", err)
    return false, err
  default:
    if count == 0 {
      log.Printf("no user with given user/password pair, user: %v\n", user.Name)
      return false, nil
    }
    return true, nil
  }
}

func (r *Repository) hasUser(ctx context.Context, name string) (bool, error) {
  var count int
  err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users WHERE " +
    "name = ?", name).Scan(&count)
  switch {
  case err != nil:
    log.Printf("query error: %v\n", err)
    return false, err
  default:
    if count == 0 {
      log.Printf("no user with name %v\n", name)
      return false, nil
    }
    return true, nil
  }
}