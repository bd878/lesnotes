package sessions

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	sessionsmodel "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Gateway struct {
	addr          string
	client        api.SessionsClient
	conn          *grpc.ClientConn
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
	g.client = api.NewSessionsClient(conn)
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) GetSession(ctx context.Context, token string) (session *sessionsmodel.Session, err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	resp, err := g.client.Get(ctx, &api.GetSessionRequest{
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	session = sessionsmodel.SessionFromProto(resp)

	return
}

func (g *Gateway) ListUserSessions(ctx context.Context, userID int64) (sessions []*sessionsmodel.Session, err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	resp, err := g.client.List(ctx, &api.ListUserSessionsRequest{
		UserId: int32(userID),
	})
	if err != nil {
		return nil, err
	}

	sessions = sessionsmodel.MapSessionsFromProto(sessionsmodel.SessionFromProto, resp.Sessions)

	return
}

func (g *Gateway) CreateSession(ctx context.Context, userID int64) (session *sessionsmodel.Session, err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	resp, err := g.client.Create(ctx, &api.CreateSessionRequest{
		UserId:         int32(userID),
	})

	session = sessionsmodel.SessionFromProto(resp)

	return
}

func (g *Gateway) RemoveSession(ctx context.Context, token string) (err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	_, err = g.client.Remove(ctx, &api.RemoveSessionRequest{
		Token:  token,
	})

	return
}

func (g *Gateway) RemoveAllUserSessions(ctx context.Context, userID int64) (err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	_, err = g.client.RemoveAll(ctx, &api.RemoveAllSessionsRequest{
		UserId: int32(userID),
	})

	return
}
