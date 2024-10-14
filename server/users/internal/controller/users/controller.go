package users

import (
  "time"
  "context"

  "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/users/internal/repository"
  "github.com/bd878/gallery/server/users/internal/controller"
)

type Repository interface {
  Add(context.Context, *model.User) error
  Has(context.Context, *model.User) (bool, error)
  Refresh(context.Context, *model.User) error
  Get(context.Context, *model.User) (*model.User, error)
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) Add(ctx context.Context, user *model.User) error {
  return c.repo.Add(ctx, user)
}

func (c *Controller) Has(ctx context.Context, user *model.User) (bool, error) {
  return c.repo.Has(ctx, user)
}

func (c *Controller) Refresh(ctx context.Context, user *model.User) error {
  return c.repo.Refresh(ctx, user)
}

func (c *Controller) Get(ctx context.Context, user *model.User) (*model.User, error) {
  var tokenExpiresTime time.Time

  result, err := c.repo.Get(ctx, &model.User{
    Token: user.Token,
    Name: user.Name,
  })
  if err == repository.ErrNoUser {
    return nil, controller.ErrNotFound
  }
  if err != nil {
    return nil, err
  }

  if result.Expires == "" {
    return nil, controller.ErrTokenExpired
  }

  err = tokenExpiresTime.UnmarshalText([]byte(result.Expires))
  if err != nil {
    return nil, err
  }

  if isTokenExpired(tokenExpiresTime) {
    return nil, controller.ErrTokenExpired
  }

  return &model.User{
    Id: result.Id,
    Name: result.Name,
    Token: result.Token,
    Expires: result.Expires,
  }, nil
}

func isTokenExpired(expiresAt time.Time) bool {
  now := time.Now()
  if expiresAt.Before(now) {
    return true
  }
  if expiresAt.Equal(now) {
    return true
  }
  return false
}
