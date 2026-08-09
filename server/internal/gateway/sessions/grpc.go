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

func (g *Gateway) GetSession(ctx context.Context, token string) (*model.Session, error) {
	resp, err := g.client.Get(ctx, &sessions.GetSessionRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return model.SessionFromProto(resp), nil
}
