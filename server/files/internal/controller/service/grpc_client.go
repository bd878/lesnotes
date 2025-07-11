package service

import (
	"context"
	"bytes"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Config struct {
	RpcAddr string
}

type Files struct {
	conf    Config
	client  api.FilesClient
	conn   *grpc.ClientConn
}

func New(cfg Config) *Files {
	f := &Files{conf: cfg}
	f.setupConnection()
	return f
}

func (f *Files) setupConnection() {
	conn, err := grpc.NewClient(f.conf.RpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*50), grpc.MaxCallSendMsgSize(1024*1024*50)))
	if err != nil {
		panic(err)
	}

	client := api.NewFilesClient(conn)

	f.conn = conn
	f.client = client
}

func (f *Files) isConnFailed() bool {
	state := f.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (f *Files) Close() {
	if f.conn != nil {
		f.conn.Close()
	}
}

type streamReader struct {
	api.Files_ReadFileStreamClient
	buf bytes.Buffer
}

func (s *streamReader) Read(p []byte) (int, error) {
	if s.buf.Available() > 0 {
		return s.buf.Read(p)
	}

	data, err := s.Recv()
	if err != nil {
		logger.Errorln("failed to receive next data")
		return 0, err
	}

	chunk, ok := data.Data.(*api.FileData_Chunk)
	if !ok {
		logger.Errorln("FileData_Chunk expected")
		return 0, errors.New("wrong format: FileData_Chunk expected")
	}

	n := copy(p, chunk.Chunk)
	if n < len(chunk.Chunk) {
		s.buf.Grow(len(chunk.Chunk))
		_, err := s.buf.Write(chunk.Chunk[n:])
		if err != nil {
			logger.Errorf("failed to write file chunks to buffer", "error", err)
			return 0, err
		}
	}

	return n, nil
}

func (f *Files) ReadFileStream(ctx context.Context, log *logger.Logger, params *model.ReadFileStreamParams) (*model.File, io.Reader, error) {
	if f.isConnFailed() {
		log.Info("conn failed, setup new connection")
		f.setupConnection()
	}

	stream, err := f.client.ReadFileStream(ctx, &api.ReadFileStreamRequest{
		Id:      params.FileID,
		UserId:  params.UserID,
	})
	if err != nil {
		log.Errorln("failed to open read stream")
		return nil, nil, err
	}

	data, err := stream.Recv()
	if err != nil {
		return nil, nil, err
	}

	meta, ok := data.Data.(*api.FileData_File)
	if !ok {
		log.Errorln("FileData_File expected")
		return nil, nil, errors.New("wrong format: FileData_File expected")
	}

	reader := &streamReader{
		Files_ReadFileStreamClient: stream,
	}

	return model.FileFromProto(meta.File), reader, nil
}

func (f *Files) SaveFileStream(ctx context.Context, log *logger.Logger, fileStream io.Reader, params *model.SaveFileParams) (
	*model.SaveFileResult, error,
) {
	if f.isConnFailed() {
		log.Info("conn failed, setup new connection")
		f.setupConnection()
	}

	stream, err := f.client.SaveFileStream(ctx)
	if err != nil {
		log.Errorw("client failed to obtain file stream", "error", err)
		return nil, err
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Name:    params.Name,
				UserId:  params.UserID,
				Private: params.Private,
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
