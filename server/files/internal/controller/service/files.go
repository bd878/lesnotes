package service

import (
	"context"
	"bytes"
	"errors"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	files "github.com/bd878/gallery/server/files/pkg/model"
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
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *streamReader) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buf.Len() > 0 {
		return s.buf.Read(p)
	}

	data, err := s.Recv()
	if err != nil {
		return 0, err
	}

	chunk, ok := data.Data.(*api.FileData_Chunk)
	if !ok {
		return 0, errors.New("wrong format: FileData_Chunk expected")
	}

	_, err = s.buf.Write(chunk.Chunk)
	if err != nil {
		return 0, err
	}

	return s.buf.Read(p)
}

func (f *Files) ReadFileStream(ctx context.Context, id int64, fileName string, public bool) (result *files.File, reader io.Reader, err error) {
	if f.isConnFailed() {
		logger.Info("conn failed, setup new connection")
		f.setupConnection()
	}

	logger.Debugw("read file stream", "id", id, "name", fileName, "public", public)

	stream, err := f.client.ReadFileStream(ctx, &api.ReadFileStreamRequest{
		Id:      id,
		Name:    fileName,
		Public:  public,
	})
	if err != nil {
		return nil, nil, err
	}

	data, err := stream.Recv()
	if err != nil {
		return nil, nil, err
	}

	meta, ok := data.Data.(*api.FileData_File)
	if !ok {
		logger.Errorln("FileData_File expected")
		return nil, nil, errors.New("wrong format: FileData_File expected")
	}

	reader = &streamReader{
		Files_ReadFileStreamClient: stream,
	}

	return files.FileFromProto(meta.File), reader, nil
}

func (f *Files) SaveFileStream(ctx context.Context, fileStream io.Reader, id, userID int64, fileName string, private bool, mime string) (err error) {
	if f.isConnFailed() {
		logger.Info("conn failed, setup new connection")
		f.setupConnection()
	}

	logger.Debugw("save file stream", "user_id", userID, "name", fileName, "private", private, "mime", mime)

	stream, err := f.client.SaveFileStream(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Id:      id,
				Name:    fileName,
				UserId:  userID,
				Private: private,
				Mime:    mime,
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

func (f *Files) ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list *files.List, err error) {
	if f.isConnFailed() {
		logger.Info("conn failed, setup new connection")
		f.setupConnection()
	}

	logger.Debugw("list files", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending, "private", private)

	resp, err := f.client.ListFiles(ctx, &api.ListFilesRequest{
		UserId:      userID,
		Limit:       limit,
		Offset:      offset,
		Asc:         ascending,
		Private:     private,
	})
	if err != nil {
		return nil, err
	}

	list = &files.List{
		Files:       files.MapFilesFromProto(files.FileFromProto, resp.Files),
		IsLastPage:  resp.IsLastPage,
		IsFirstPage: offset == 0,
		Count:       int32(len(resp.Files)),
		// TODO: total
		Offset:      offset,
	}

	return
}