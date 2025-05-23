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
	DeleteUser(ctx context.Context, log *logger.Logger, params *model.DeleteUserParams) error
	RefreshToken(ctx context.Context, log *logger.Logger, user *model.User) error
	DeleteToken(ctx context.Context,  log *logger.Logger, params *model.DeleteTokenParams) error
	GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error)
}

type Controller struct {
	repo Repository
}

func New(repo Repository) *Controller {
	return &Controller{repo}
}

func (c *Controller) AddUser(ctx context.Context, log *logger.Logger, params *model.AddUserParams) (int32, error) {
	params.User.ID = utils.RandomID()
	err := c.repo.AddUser(ctx, log, params.User)
	return params.User.ID, err
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

func (c *Controller) DeleteUser(ctx context.Context, log *logger.Logger, params *model.DeleteUserParams) error {
	return c.repo.DeleteUser(ctx, log, params)
}

func (c *Controller) GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error) {
	log.Debugw("controller get user", "id", params.ID, "token", params.Token)
	user, err := c.repo.GetUser(ctx, log, params)
	if err != nil {
		log.Errorln("failed to get user")
		return nil, err
	}

	if user.ExpiresUTCNano == 0 {
		return nil, controller.ErrTokenExpired
	}

	if isTokenExpired(user.ExpiresUTCNano) {
		return nil, controller.ErrTokenExpired
	}

	log.Debugw("controller get user ok", "id", params.ID, "token", params.Token)
	return user, nil
}

func isTokenExpired(expiresUtcNano int64) bool {
	return time.Now().UnixNano() > expiresUtcNano
}
