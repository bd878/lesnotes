package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/db/users/pkg/loadbalance"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Gateway struct {
	addr   string
	client users.UsersClient
	conn   *grpc.ClientConn
}

func New(addr string) *Gateway {
	g := &Gateway{addr: addr}
	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.addr,
		), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	g.conn = conn
	g.client = users.NewUsersClient(conn)
	return nil
	return g
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
