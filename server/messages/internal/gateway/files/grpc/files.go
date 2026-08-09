package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api/files"
	"github.com/bd878/gallery/server/internal/rpc"
)

type Gateway struct {
	filesAddr string
	client    files.FilesClient
	conn      *grpc.ClientConn
}

func New(filesAddr string) *Gateway {
	g := &Gateway{filesAddr: filesAddr}

	conn, err := rpc.NewClient(
		g.filesAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = files.NewFilesClient(conn)

	return g
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
