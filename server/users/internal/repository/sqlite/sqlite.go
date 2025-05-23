package repository

import (
	"errors"
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Repository struct {
	db          *sql.DB
	resetTokenStmt *sql.Stmt
	selectByIDStmt *sql.Stmt
}

func New(dbPath string) *Repository {
	db, err := sql.Open("sqlite3", "file:" + dbPath)
	if err != nil {
		panic(err)
	}

	resetTokenStmt := utils.Must(db.Prepare(`
UPDATE users SET
	token = "",
	expires_utc_nano = 0
WHERE token = :token
;`,
	))

	selectByIDStmt := utils.Must(db.Prepare(`
SELECT id, name, password, token, expires_utc_nano
FROM users
WHERE id = :id
;`,
	))

	return &Repository{
		db: db,
		resetTokenStmt: resetTokenStmt,
		selectByIDStmt: selectByIDStmt,
	}
}

func (r *Repository) AddUser(ctx context.Context, log *logger.Logger, user *model.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users(id,name,password,token,expires_utc_nano)" +
		"VALUES(?,?,?,?,?)", user.ID, user.Name, user.Password, user.Token, user.ExpiresUTCNano)
	if err != nil {
		log.Errorf("query error: %v\n", err)
	}
	return err
}

func (r *Repository) HasUser(ctx context.Context, log *logger.Logger, user *model.User) (bool, error) {
	if user.Password == "" {
		return r.hasUser(ctx, log, user.Name)
	}

	return r.hasUserAndPassword(ctx, log, user)
}

func (r *Repository) GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error) {
	if params.Token != "" {
		return r.getByToken(ctx, log, params.Token)
	} else if params.Name != "" {
		return r.getByUserName(ctx, log, params.Name)
	} else if params.ID != 0 {
		return r.getByID(ctx, log, params.ID)
	}
	return nil, errors.New("not implemented")
}

func (r *Repository) RefreshToken(ctx context.Context, log *logger.Logger, user *model.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET token = ?, expires_utc_nano = ? WHERE name = ?",
		user.Token, user.ExpiresUTCNano, user.Name)
	return err
}

func (r *Repository) DeleteToken(ctx context.Context, log *logger.Logger, params *model.DeleteTokenParams) error {
	_, err := r.resetTokenStmt.ExecContext(ctx, sql.Named("token", params.Token))
	if err != nil {
		log.Error("cannot delete token")
	}
	return err
}

func (r *Repository) getByID(ctx context.Context, log *logger.Logger, id int32) (*model.User, error) {
	var (
		_id int32
		expiresUtcNano int64
		passwordCol sql.NullString
		tokenCol sql.NullString
		nameCol sql.NullString
	)

	err := r.selectByIDStmt.QueryRowContext(ctx, sql.Named("id", id)).Scan(
		&_id, &nameCol, &passwordCol, &tokenCol, &expiresUtcNano)
	if err != nil {
		log.Errorw("repository failed to find user by id, sqlite error", "user_id", id)
		return nil, err
	}

	user := &model.User{
		ID: id,
		ExpiresUTCNano: expiresUtcNano,
	}

	if passwordCol.Valid {
		user.Password = passwordCol.String
	}

	if tokenCol.Valid {
		user.Token = tokenCol.String
	}

	if nameCol.Valid {
		user.Name = nameCol.String
	}

	return user, nil
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
		log.Errorf("no rows for name %v\n", name)
		return nil, errors.New("no user")

	case err != nil:
		log.Errorf("query error: %v\n", err)
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
		log.Errorf("no rows for token %v\n", token)
		return nil, errors.New("no user")

	case err != nil:
		log.Errorf("query error: %v\n", err)
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
		log.Errorf("query error: %v\n", err)
		return false, err
	default:
		if count == 0 {
			log.Errorf("no user with given user/password pair, user: %v\n", user.Name)
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
		log.Errorf("query error: %v\n", err)
		return false, err
	default:
		if count == 0 {
			log.Errorf("no user with name %v\n", name)
			return false, nil
		}
		return true, nil
	}
}