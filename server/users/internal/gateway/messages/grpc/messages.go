package messages

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
)

type Gateway struct {
	messagesAddr string
	client    api.MessagesClient
	conn      *grpc.ClientConn
}

func New(messagesAddr string) *Gateway {
	g := &Gateway{messagesAddr: messagesAddr}
	g.setupConnection()
	return g
}

func (g *Gateway) setupConnection() {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.messagesAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = api.NewMessagesClient(conn)
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) DeleteAllUserMessages(ctx context.Context, userID int64) (err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	_, err = g.client.DeleteAllUserMessages(ctx, &api.DeleteAllUserMessagesRequest{
		UserId: userID,
	})
	return
}