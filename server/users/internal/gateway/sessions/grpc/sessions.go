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
	addr          string
	client        api.SessionsClient
	conn          *grpc.ClientConn
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
		logger.Debugw("sessions connection failed", "state", state)
		return true
	}
	return false
}

func (g *Gateway) GetSession(ctx context.Context, token string) (session *model.Session, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	resp, err := g.client.Get(ctx, &api.GetSessionRequest{
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	session = model.SessionFromProto(resp)

	return
}

func (g *Gateway) ListUserSessions(ctx context.Context, userID int64) (sessions []*model.Session, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	resp, err := g.client.List(ctx, &api.ListUserSessionsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	sessions = model.MapSessionsFromProto(model.SessionFromProto, resp.Sessions)

	return
}

func (g *Gateway) CreateSession(ctx context.Context, userID int64) (session *model.Session, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	resp, err := g.client.Create(ctx, &api.CreateSessionRequest{
		UserId:         userID,
	})
	if err != nil {
		return nil, err
	}

	session = model.SessionFromProto(resp)

	return
}

func (g *Gateway) RemoveSession(ctx context.Context, token string) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	_, err = g.client.Remove(ctx, &api.RemoveSessionRequest{
		Token:  token,
	})

	return
}

func (g *Gateway) RemoveAllUserSessions(ctx context.Context, userID int64) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	_, err = g.client.RemoveAll(ctx, &api.RemoveAllSessionsRequest{
		UserId: userID,
	})

	return
}
