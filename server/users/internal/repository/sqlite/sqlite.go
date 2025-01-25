package repository

import (
  "errors"
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/users/pkg/model"
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

func (r *Repository) AddUser(ctx context.Context, log *logger.Logger, user *model.User) error {
  _, err := r.db.ExecContext(ctx, "INSERT INTO users(id,name,password,token,expires_utc_nano)" +
    "VALUES(?,?,?,?,?)", user.ID, user.Name, user.Password, user.Token, user.ExpiresUTCNano)
  if err != nil {
    log.Error("query error: %v\n", err)
  }
  return err
}

func (r *Repository) HasUser(ctx context.Context, log *logger.Logger, user *model.User) (bool, error) {
  if user.Password == "" {
    return r.hasUser(ctx, log, user.Name)
  }

  return r.hasUserAndPassword(ctx, log, user)
}

func (r *Repository) GetUser(ctx context.Context, log *logger.Logger, user *model.User) (*model.User, error) {
  if user.Token != "" {
    return r.getByToken(ctx, log, user.Token)
  } else if user.Name != "" {
    return r.getByUserName(ctx, log, user.Name)
  }
  return nil, errors.New("not implemented")
}

func (r *Repository) RefreshToken(ctx context.Context, log *logger.Logger, user *model.User) error {
  _, err := r.db.ExecContext(ctx, "UPDATE users SET token = ?, expires_utc_nano = ? WHERE name = ?",
    user.Token, user.ExpiresUTCNano, user.Name)
  return err
}

func (r *Repository) getByUserName(ctx context.Context, log *logger.Logger, name string) (*model.User, error) {
  var password, token string
  var expiresUtcNano int64
  var id int32

  err := r.db.QueryRowContext(ctx, "SELECT id, name, password, token, expires_utc_nano FROM users WHERE " +
    "name = ?", name).Scan(&id, &name, &password, &token, &expiresUtcNano)

  msg := &model.User{
    ID:                id,
    Name:              name,
    Password:          password,
    Token:             token,
    ExpiresUTCNano:    expiresUtcNano,
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

func (r *Repository) getByToken(ctx context.Context, log *logger.Logger, token string) (*model.User, error) {
  var name, password string
  var expiresUtcNano int64
  var id int32

  err := r.db.QueryRowContext(ctx, "SELECT id, name, password, token, expires_utc_nano FROM users WHERE " +
    "token = ?", token).Scan(&id, &name, &password, &token, &expiresUtcNano)
  switch {
  case err == sql.ErrNoRows:
    log.Error("no rows for token %v\n", token)
    return nil, errors.New("no user")

  case err != nil:
    log.Error("query error: %v\n", err)
    return nil, err

  default:
    return &model.User{
      ID:                     id,
      Name:                   name,
      Password:               password,
      Token:                  token,
      ExpiresUTCNano:         expiresUtcNano,
    }, nil
  }
}

func (r *Repository) hasUserAndPassword(ctx context.Context, log *logger.Logger, user *model.User) (bool, error) {
  var count int
  err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users WHERE " +
    "name = ? AND password = ?", user.Name, user.Password).Scan(&count)
  switch {
  case err != nil:
    log.Error("query error: %v\n", err)
    return false, err
  default:
    if count == 0 {
      log.Error("no user with given user/password pair, user: %v\n", user.Name)
      return false, nil
    }
    return true, nil
  }
}

func (r *Repository) hasUser(ctx context.Context, log *logger.Logger, name string) (bool, error) {
  var count int
  err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users WHERE " +
    "name = ?", name).Scan(&count)
  switch {
  case err != nil:
    log.Error("query error: %v\n", err)
    return false, err
  default:
    if count == 0 {
      log.Error("no user with name %v\n", name)
      return false, nil
    }
    return true, nil
  }
}