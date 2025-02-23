package users

import (
  "time"
  "context"

  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/users/internal/controller"
)

type Repository interface {
  AddUser(ctx context.Context, log *logger.Logger, user *model.User) error
  HasUser(ctx context.Context, log *logger.Logger, user *model.User) (bool, error)
  RefreshToken(ctx context.Context, log *logger.Logger, user *model.User) error
  DeleteToken(ctx context.Context,  log *logger.Logger, params *model.DeleteTokenParams) error
  GetUser(ctx context.Context, log *logger.Logger, user *model.User) (*model.User, error)
}

type Controller struct {
  repo Repository
}

func New(repo Repository) *Controller {
  return &Controller{repo}
}

func (c *Controller) AddUser(ctx context.Context, log *logger.Logger, params *model.AddUserParams) error {
  params.User.ID = utils.RandomID()
  return c.repo.AddUser(ctx, log, params.User)
}

func (c *Controller) HasUser(ctx context.Context, log *logger.Logger, params *model.HasUserParams) (bool, error) {
  return c.repo.HasUser(ctx, log, params.User)
}

func (c *Controller) RefreshToken(ctx context.Context, log *logger.Logger, params *model.RefreshTokenParams) error {
  return c.repo.RefreshToken(ctx, log, params.User)
}

func (c *Controller) DeleteToken(ctx context.Context, log *logger.Logger, params *model.DeleteTokenParams) error {
  return c.repo.DeleteToken(ctx, log, params)
}

func (c *Controller) GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error) {
  result, err := c.repo.GetUser(ctx, log, &model.User{
    Token: params.User.Token,
    Name: params.User.Name,
  })
  if err != nil {
    log.Errorln("failed to get user")
    return nil, err
  }

  if result.ExpiresUTCNano == 0 {
    return nil, controller.ErrTokenExpired
  }

  if isTokenExpired(result.ExpiresUTCNano) {
    return nil, controller.ErrTokenExpired
  }

  return &model.User{
    ID:              result.ID,
    Name:            result.Name,
    Token:           result.Token,
    ExpiresUTCNano:  result.ExpiresUTCNano,
  }, nil
}

func isTokenExpired(expiresUtcNano int64) bool {
  return time.Now().UnixNano() > expiresUtcNano
}
