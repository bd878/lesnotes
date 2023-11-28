package users

import (
  "context"

  "github.com/bd878/gallery/server/user/pkg/model"
  "github.com/bd878/gallery/server/user/internal/repository"
  "github.com/bd878/gallery/server/user/internal/controller"
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
  result, err := c.repo.Get(ctx, &model.User{
    Token: user.Token,
  })
  if err == repository.ErrNoUser {
    return nil, controller.ErrNotFound
  }
  if err != nil {
    return nil, err
  }
  // TODO: check for expire
  if result.Token == user.Token {
    return &model.User{
      Id: result.Id,
      Name: result.Name,
      Token: result.Token,
      Expires: result.Expires,
    }, nil
  }
  return nil, controller.ErrTokenInvalid
}
