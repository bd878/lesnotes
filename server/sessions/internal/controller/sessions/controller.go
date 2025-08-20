package sessions

import (
	"context"
	"time"

	"github.com/bd878/gallery/server/sessions/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

type Repository interface {
	Save(ctx context.Context, userID int32, token string, expiresUTCNano int64) (err error)
	Get(ctx context.Context, token string) (session *model.Session, err error)
	List(ctx context.Context, userID int32) (sessions []*model.Session, err error)
	Delete(ctx context.Context, token string) (err error)
	DeleteAll(ctx context.Context, userID int32) (err error)
}

type Controller struct {
	repo Repository
}

func New(repo Repository) *Controller {
	return &Controller{repo}
}

func (c *Controller) CreateSession(ctx context.Context, userID int32) (session *model.Session, err error) {
	token := utils.RandomString(10)
	expiresUTCNano := time.Now().Add(time.Hour * 24 * 5).UnixNano()

	err = c.repo.Save(ctx, userID, token, expiresUTCNano)
	if err != nil {
		return
	}

	session = &model.Session{
		UserID:         userID,
		Token:          token,
		ExpiresUTCNano: expiresUTCNano,
	}

	return
}

func (c *Controller) GetSession(ctx context.Context, token string) (*model.Session, error) {
	return c.repo.Get(ctx, token)
}

func (c *Controller) ListUserSessions(ctx context.Context, userID int32) ([]*model.Session, error) {
	return c.repo.List(ctx, userID)
}

func (c *Controller) RemoveSession(ctx context.Context, token string) error {
	return c.repo.Delete(ctx, token)
}

func (c *Controller) RemoveUserSessions(ctx context.Context, userID int32) error {
	return c.repo.DeleteAll(ctx, userID)
}