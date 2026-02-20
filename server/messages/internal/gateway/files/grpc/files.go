package grpc

import (
	"io"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	files "github.com/bd878/gallery/server/files/pkg/model"
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
	conn, err := grpc.NewClient(g.filesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (g *Gateway) ReadBatchFiles(ctx context.Context, fileIDs []int64, userID int64) (result map[int64]*files.File, err error) {
	if g.isConnFailed() {
		logger.Info("conn failed, setup new connection")
		g.setupConnection()
	}

	logger.Debugw("read batch files", "file_ids", fileIDs, "user_id", userID)

	batch, err := g.client.ReadBatchFiles(ctx, &api.ReadBatchFilesRequest{
		UserId: userID,
		Ids:    fileIDs,
	})
	if err != nil {
		return nil, err
	}

	result = files.MapFilesDictFromProto(files.FileFromProto, batch.Files)

	return
}

func (g *Gateway) ReadFile(ctx context.Context, userID, fileID int64) (resp *files.File, err error) {
	logger.Debugw("read file", "user_id", userID, "file_id", fileID)

	file, err := g.client.ReadFile(ctx, &api.ReadFileRequest{
		UserId: userID,
		Id:     fileID,
	})
	if err != nil {
		return nil, err
	}

	resp = files.FileFromProto(file)

	return
}

// copied from files/internal/controller/service
func (g *Gateway) SaveFile(ctx context.Context, fileStream io.Reader, id, userID int64, fileName string, private bool, mime string) (err error) {
	if g.isConnFailed() {
		g.setupConnection()
	}

	logger.Debugw("save file", "id", id, "user_id", userID, "file_name", fileName, "private", private, "mime", mime)

	stream, err := g.client.SaveFileStream(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Id:      id,
				Private: private,
				Mime:    mime,
				Name:    fileName,
				UserId:  userID,
			},
		},
	})
	if err != nil {
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := fileStream.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&api.FileData{
			Data: &api.FileData_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()

	return
}