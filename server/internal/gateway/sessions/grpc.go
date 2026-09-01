package sessions

import (
	"context"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/db/sessions/pkg/model"
)

type Gateway struct {
	client sessions.SessionsClient
}

func New(container di.Container) *Gateway {
	client := container.Get("sessionsClient").(sessions.SessionsClient)

	return &Gateway{client}
}

func (g *Gateway) GetSession(ctx context.Context, token string) (*model.Session, error) {
	resp, err := g.client.Get(ctx, &sessions.GetSessionRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return model.SessionFromProto(resp), nil
}
