package users

import (
  "context"

  "github.com/bd878/gallery/server/pkg/model"
)

type Repository interface {
  AddUser(context.Context, *model.User) error
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) Add(ctx context.Context, usr *model.User) error {
  err := c.repo.AddUser(ctx, usr)
  if err != nil {
    return err
  }
  return nil
}
