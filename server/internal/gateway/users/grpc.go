package users

import (
	"context"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Gateway struct {
	client users.UsersClient
}

func New(container di.Container) *Gateway {
	client := container.Get("usersClient").(users.UsersClient)

	return &Gateway{client}
}

func (g *Gateway) GetUser(ctx context.Context, userID int64) (*model.User, error) {

	resp, err := g.client.GetUser(ctx, &users.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return model.UserFromProto(resp), nil
}
