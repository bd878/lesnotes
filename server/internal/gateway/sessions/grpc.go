package sessions

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/sessions/pkg/loadbalance"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

type Gateway struct {
	addr            string
	client          api.SessionsClient
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
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = api.NewSessionsClient(conn)

	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugw("gateway conn failed", "state", state)
		return true
	}
	return false
}

func (g *Gateway) GetSession(ctx context.Context, token string) (*model.Session, error) {
	if g.isConnFailed() {
		if err := g.setupConnection(); err != nil {
			return nil, err
		}
	}
	resp, err := g.client.Get(ctx, &api.GetSessionRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return model.SessionFromProto(resp), nil
}
