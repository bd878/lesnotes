package users

import (
  "context"

  "github.com/bd878/gallery/server/user/pkg/model"
)

type Repository interface {
  Add(context.Context, *model.User) error
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) Add(ctx context.Context, usr *model.User) error {
  err := c.repo.Add(ctx, usr)
  if err != nil {
    return err
  }
  return nil
}
