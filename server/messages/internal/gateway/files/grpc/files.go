package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	"github.com/bd878/gallery/server/internal/grpcutil"
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

func (g *Gateway) setupConnection() {
	conn, err := grpcutil.ServiceConnection(context.Background(), g.filesAddr)
	if err != nil {
		panic(err)
	}

	g.conn = conn
	g.client = api.NewFilesClient(conn)
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (
	*model.ReadBatchFilesResult, error,
) {
	if g.isConnFailed() {
		log.Info("conn failed, setup new connection")
		g.setupConnection()
	}

	batch, err := g.client.ReadBatchFiles(ctx, &api.ReadBatchFilesRequest{
		UserId: params.UserID,
		Ids:    params.IDs,
	})
	if err != nil {
		log.Errorw("client failed to read batch files", "error", err)
		return nil, err
	}

	return &model.ReadBatchFilesResult{
		Files: filesmodel.MapFilesDictFromProto(filesmodel.FileFromProto, batch.Files),
	}, nil
}

func (g *Gateway) ReadFile(ctx context.Context, log *logger.Logger, userID, fileID int32) (
	*filesmodel.File, error,
) {
	file, err := g.client.ReadFile(ctx, &api.ReadFileRequest{
		UserId: userID,
		Id: fileID,
	})
	if err != nil {
		log.Errorw("client failed to read one file", "user_id", userID, "file_id", fileID)
		return nil, err
	}

	return filesmodel.FileFromProto(file), nil
}