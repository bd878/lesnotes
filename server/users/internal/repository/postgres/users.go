package postgres

import (
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/users/pkg/model"
)

type UsersRepository struct {
	usersTableName      string
	premiumsTableName   string
	pool                *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool, usersTableName, premiumsTableName string) *UsersRepository {
	return &UsersRepository{
		usersTableName:    usersTableName,
		premiumsTableName: premiumsTableName,
		pool:              pool,
	}
}

func (r *UsersRepository) Save(ctx context.Context, id int64, login, salt string, metadata []byte, createdAt, updatedAt string) (err error) {
	const query = "INSERT INTO %s(id, login, salt, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = r.pool.Exec(ctx, r.usersTable(query), id, login, salt, metadata, createdAt, updatedAt)

	return
}

func (r *UsersRepository) Delete(ctx context.Context, id int64) (err error) {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.usersTable(query), id)

	return
}

func (r *UsersRepository) FindByID(ctx context.Context, id int64) (user *model.User, err error) {
	query := fmt.Sprintf(`
SELECT DISTINCT
	u.login,
	u.salt,
	u.metadata,
	u.created_at,
	u.updated_at,
	coalesce(p.is_premium, false) as is_premium
FROM %s u
LEFT JOIN (
	SELECT id, EXTRACT(EPOCH FROM expires_at - NOW()) > 0 as is_premium
	FROM %s
	WHERE id = $1
) p
ON u.id = p.id
WHERE u.id = $1
`, r.usersTableName, r.premiumsTableName)

	user = &model.User{
		ID:     id,
	}

	var createdAt, updatedAt *time.Time

	err = r.pool.QueryRow(ctx, query, id).Scan(&user.Login, &user.HashedPassword,
		&user.Metadata, &createdAt, &updatedAt, &user.IsPremium)
	if err != nil {
		return
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *UsersRepository) FindByLogin(ctx context.Context, login string) (user *model.User, err error) {
	query := fmt.Sprintf(`
SELECT DISTINCT
	u.id,
	u.salt,
	u.metadata,
	u.created_at,
	u.updated_at,
	coalesce(p.is_premium, false) as is_premium
FROM %s u
LEFT JOIN (
	SELECT id, EXTRACT(EPOCH FROM expires_at - NOW()) > 0 as is_premium
	FROM %s
) p
ON u.id = p.id
WHERE u.login = $1
`, r.usersTableName, r.premiumsTableName)

	user = &model.User{
		Login:   login,
	}

	var createdAt, updatedAt *time.Time

	err = r.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.HashedPassword,
		&user.Metadata, &createdAt, &updatedAt, &user.IsPremium)
	if err != nil {
		return
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	return
}

func (r *UsersRepository) Update(ctx context.Context, id int64, login *string, metadata []byte, updatedAt string) (err error) {
	const query = "UPDATE %s SET login = $2, metadata = $3, updated_at = $4 WHERE id = $1"

	_, err = r.pool.Exec(ctx, r.usersTable(query), id, login, metadata, updatedAt)

	return
}

func (r *UsersRepository) MakePremium(ctx context.Context, id int64, invoiceID, createdAt, expiresAt string) (err error) {
	const query = "INSERT INTO %s(id, invoice_id, created_at, expires_at) VALUES ($1, $2, $3, $4)"

	_, err = r.pool.Exec(ctx, r.premiumsTable(query), id, invoiceID, createdAt, expiresAt)

	return
}

func (r UsersRepository) usersTable(query string) string {
	return fmt.Sprintf(query, r.usersTableName)
}

func (r UsersRepository) premiumsTable(query string) string {
	return fmt.Sprintf(query, r.premiumsTableName)
}