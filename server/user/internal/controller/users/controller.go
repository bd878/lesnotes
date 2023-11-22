package users

import (
  "context"

  "github.com/bd878/gallery/server/user/pkg/model"
)

type Repository interface {
  Add(context.Context, *model.User) error
  Has(context.Context, *model.User) (bool, error)
  Refresh(context.Context, *model.User) error
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