package messages

import (
  "context"

  "github.com/bd878/gallery/server/pkg/model"
)

type Repository interface {
  Put(context.Context, string) error
  GetAll(context.Context) ([]model.Message, error)
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) SaveMessage(ctx context.Context, msg string) error {
  err := c.repo.Put(ctx, msg)
  if err != nil {
    return err
  }
  return nil
}

func (c *Controller) ReadAllMessages(ctx context.Context) ([]model.Message, error) {
  v, err := c.repo.GetAll(ctx)
  if err != nil {
    return nil, err
  }
  return v, nil
}
