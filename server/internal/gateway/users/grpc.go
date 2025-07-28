package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Gateway struct {
	addr            string
	client          api.UsersClient
	conn            *grpc.ClientConn
}

func New(addr string) *Gateway {
	g := &Gateway{addr: addr}
	g.setupConnection()
	return g
}

func (g *Gateway) setupConnection() {
	conn, err := grpc.NewClient(g.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = api.NewUsersClient(conn)
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) GetUser(ctx context.Context, userID int32) (*usersmodel.User, error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	resp, err := g.client.GetUser(ctx, &api.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return usersmodel.UserFromProto(resp), nil
}
