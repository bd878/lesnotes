package messages

import (
  "time"
  "context"

  usermodel "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository interface {
  Put(context.Context, *model.Message) error
  Get(context.Context, usermodel.UserId) ([]model.Message, error)
}

type Controller struct {
  repo Repository
}

func New(repo Repository) (*Controller, error) {
  return &Controller{repo}, nil
}

func (c *Controller) SaveMessage(ctx context.Context, msg *model.Message) error {
  msg.CreateTime = time.Now().String()
  err := c.repo.Put(ctx, msg)
  if err != nil {
    return err
  }
  return nil
}

func (c *Controller) ReadUserMessages(ctx context.Context, userId usermodel.UserId) ([]model.Message, error) {
  v, err := c.repo.Get(ctx, userId)
  if err != nil {
    return nil, err
  }
  return v, nil
}
