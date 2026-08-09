package sessions

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"
	"github.com/bd878/gallery/server/db/sessions/pkg/model"
	"github.com/bd878/gallery/server/internal/rpc"
)

type Gateway struct {
	addr   string
	client sessions.SessionsClient
	conn   *grpc.ClientConn
}

func New(addr string) *Gateway {
	g := &Gateway{addr: addr}

	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.addr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = sessions.NewSessionsClient(conn)

	return g
}

func (g *Gateway) GetSession(ctx context.Context, token string) (session *model.Session, err error) {

	resp, err := g.client.Get(ctx, &sessions.GetSessionRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	session = model.SessionFromProto(resp)

	return
}

func (g *Gateway) ListUserSessions(ctx context.Context, userID int64) (list []*model.Session, err error) {

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

	resp, err := g.client.Create(ctx, &sessions.CreateSessionRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	session = model.SessionFromProto(resp)

	return
}

func (g *Gateway) RemoveSession(ctx context.Context, token string) (err error) {

	_, err = g.client.Remove(ctx, &sessions.RemoveSessionRequest{
		Token: token,
	})

	return
}

func (g *Gateway) RemoveAllUserSessions(ctx context.Context, userID int64) (err error) {

	_, err = g.client.RemoveAll(ctx, &sessions.RemoveAllSessionsRequest{
		UserId: userID,
	})

	return
}
