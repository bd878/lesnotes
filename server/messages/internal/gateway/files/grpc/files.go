package grpc

import (
	"context"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/api/files"
)

type Gateway struct {
	client    files.FilesClient
}

func New(container di.Container) *Gateway {
	client := container.Get("filesClient").(files.FilesClient)

	return &Gateway{client}
}

func (g *Gateway) ReadMessageFiles(ctx context.Context, messageID int64, userIDs []int64) (list []*files.File, err error) {
	resp, err := g.client.ReadMessageFiles(ctx, &files.ReadMessageFilesRequest{
		Id:      messageID,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	list = resp.Files

	return
}
