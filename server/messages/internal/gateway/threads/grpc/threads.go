package grpc

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/threads/pkg/loadbalance"
)

type Gateway struct {
	addr    string
	client  api.ThreadsClient
	conn    *grpc.ClientConn
}

func New(addr string) *Gateway {
	gateway := &Gateway{addr: addr}

	gateway.setupConnection()

	return gateway
}

func (g *Gateway) setupConnection() error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.addr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = api.NewThreadsClient(conn)

	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) CreateThread(ctx context.Context, id, userID, parentID int64, name string, private bool) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	_, err = g.client.Create(ctx, &api.CreateRequest{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		Name:      name,
		Private:   private,
	})

	return
}

func (g *Gateway) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	_, err = g.client.Delete(ctx, &api.DeleteRequest{
		Id:     id,
		UserId: userID,
	})

	return
}

func (g *Gateway) UpdateThread(ctx context.Context, id, userID, parentID int64) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	_, err = g.client.Update(ctx, &api.UpdateRequest{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		Private:   -1,
	})

	return
}
