package grpc

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/loadbalance"
	"github.com/bd878/gallery/server/users/pkg/model"
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

func (g *Gateway) setupConnection() error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.addr,
		), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = api.NewUsersClient(conn)

	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugw("users gateway conn failed", "state", state.String())
		return true
	}
	return false
}

func (g *Gateway) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	if g.isConnFailed() {
		if err := g.setupConnection(); err != nil {
			return nil, err
		}
	}

	resp, err := g.client.GetUser(ctx, &api.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return model.UserFromProto(resp), nil
}
