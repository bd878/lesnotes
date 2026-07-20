package sessions

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"
	"github.com/bd878/gallery/server/db/sessions/pkg/model"
)

type Gateway struct {
	addr          string
	client        sessions.SessionsClient
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
	g.client = sessions.NewSessionsClient(conn)

	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		logger.Debugw("sessions conn failed", "state", state.String())
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

	resp, err := g.client.Get(ctx, &sessions.GetSessionRequest{
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	session = model.SessionFromProto(resp)

	return
}

func (g *Gateway) ListUserSessions(ctx context.Context, userID int64) (list []*model.Session, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	resp, err := g.client.List(ctx, &sessions.ListUserSessionsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	list = model.MapSessionsFromProto(model.SessionFromProto, resp.Sessions)

	return
}

func (g *Gateway) CreateSession(ctx context.Context, userID int64) (session *model.Session, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	resp, err := g.client.Create(ctx, &sessions.CreateSessionRequest{
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

	_, err = g.client.Remove(ctx, &sessions.RemoveSessionRequest{
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

	_, err = g.client.RemoveAll(ctx, &sessions.RemoveAllSessionsRequest{
		UserId: userID,
	})

	return
}
