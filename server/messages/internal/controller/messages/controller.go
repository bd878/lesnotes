package messages

import (
  "time"
  "context"

  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository interface {
  Put(context.Context, *model.Message) error
  GetAll(context.Context) ([]model.Message, error)
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) SaveMessage(ctx context.Context, msg *model.Message) error {
  msg.CreateTime = time.Now().String()
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
