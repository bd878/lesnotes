package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Gateway struct {
	filesAddr string
	client    api.FilesClient
	conn      *grpc.ClientConn
}

func New(filesAddr string) *Gateway {
	g := &Gateway{filesAddr: filesAddr}
	g.setupConnection()
	return g
}

func (g *Gateway) setupConnection() (err error) {
	conn, err := grpc.NewClient(g.filesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = api.NewFilesClient(conn)
	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugw("files conn failed", "state", state.String())
		return true
	}
	return false
}

func (g *Gateway) ReadMessageFiles(ctx context.Context, messageID int64, userIDs []int64) (list []*api.File, err error) {
	if g.isConnFailed() {
		logger.Info("conn failed, setup new connection")
		if err := g.setupConnection(); err != nil {
			return nil, err
		}
	}

	resp, err := g.client.ReadMessageFiles(ctx, &api.ReadMessageFilesRequest{
		Id: messageID,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	list = resp.Files

	return
}
