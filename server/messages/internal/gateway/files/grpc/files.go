package grpc

import (
	"context"
	"io"

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

// copied from files/internal/controller/service
func (g *Gateway) SaveFile(ctx context.Context, log *logger.Logger, fileStream io.Reader, params *model.SaveFileParams) (
	*model.SaveFileResult, error,
) {
	if g.isConnFailed() {
		log.Info("conn failed, setup new connection")
		g.setupConnection()
	}

	stream, err := g.client.SaveFileStream(ctx)
	if err != nil {
		log.Errorw("client failed to obtain file stream", "error", err)
		return nil, err
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Name:    params.Name,
				UserId:  params.UserID,
			},
		},
	})
	if err != nil {
		log.Errorw("failed to save file meta", "error", err)
		return nil, err
	}

	buffer := make([]byte, 1024)
	for {
		n, err := fileStream.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorw("failed to read file data in buffer", "error", err)
			return nil, err
		}

		err = stream.Send(&api.FileData{
			Data: &api.FileData_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			log.Errorw("failed to send chunk on file server", "error", err)
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorw("failed to close and recv result", "error", err)
		return nil, err
	}

	return &model.SaveFileResult{
		ID:              res.File.Id,
		CreateUTCNano:   res.File.CreateUtcNano,
	}, nil
}