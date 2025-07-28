package sessions

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	sessionsmodel "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Gateway struct {
	addr            string
	client          api.SessionsClient
	conn            *grpc.ClientConn
}

func New(addr string) *Gateway {
	g := &Gateway{addr: addr}
	g.setupConnection()
	return &Gateway{conn}
}

func (g *Gateway) setupConnection() {
	conn, err := grpc.NewClient(g.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = api.NewSessionsClient(conn)
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) RestoreSession(ctx context.Context, token string) (*sessionsmodel.Session, error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	resp, err := g.client.Get(ctx, &api.GetSessionRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return sessionsmodel.SessionFromProto(resp), nil
}
