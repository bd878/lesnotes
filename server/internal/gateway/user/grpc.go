package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Gateway struct {
	userAddr string
}

func New(userAddr string) *Gateway {
	return &Gateway{userAddr}
}

func (g *Gateway) Auth(ctx context.Context, log *logger.Logger, token string) (*usermodel.User, error) {
	conn, err := grpc.NewClient(g.userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorw("failed to establish connection with user service", "error", err)
		return nil, err
	}
	defer conn.Close()
	client := api.NewUsersClient(conn)
	resp, err := client.Auth(ctx, &api.AuthUserRequest{Token: token})
	if err != nil {
		log.Errorw("failed to authenticate on client", "error", err)
		return nil, err
	}
	return usermodel.UserFromProto(resp.User), nil
}

func (g *Gateway) GetPublicUser(ctx context.Context, log *logger.Logger, id int32) (*usermodel.User, error) {
	conn, err := grpc.NewClient(g.userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("gateway failed to establish connection with user service")
		return nil, err
	}
	defer conn.Close()
	client := api.NewUsersClient(conn)
	user, err := client.GetUser(ctx, &api.GetUserRequest{
		SearchKey: &api.GetUserRequest_Id{
			Id: id,
		},
	})
	if err != nil {
		log.Errorw("gateway failed to get public user", "id", id)
		return nil, err
	}
	return usermodel.UserFromProto(user), nil
}