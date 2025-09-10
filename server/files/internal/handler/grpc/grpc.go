package grpc

import (
	"io"
	"context"
	"bytes"
	"sync"
	"errors"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Repository interface {
	SaveFile(ctx context.Context, reader io.Reader, id, userID int64, private bool, name, mime string) (err error)
	GetMeta(ctx context.Context, id int64, fileName string) (file *model.File, err error)
	DeleteFile(ctx context.Context, ownerID, id int64) (err error)
	ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error)
}

type Handler struct {
	api.UnimplementedFilesServer
	repo       Repository
}

func New(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ReadBatchFiles(ctx context.Context, req *api.ReadBatchFilesRequest) (*api.ReadBatchFilesResponse, error) {
	logger.Debugw("read batch files", "user_id", req.UserId, "ids", req.Ids)

	files := make(map[int64]*model.File, len(req.Ids))
	for _, id := range req.Ids {
		files[id] = &model.File{
			ID:     id,
			UserID: req.UserId,
		}

		file, err := h.repo.GetMeta(ctx, id, "")
		if err != nil {
			files[id].Error = "can not find file"
			logger.Errorw("failed to read file", "user_id", req.UserId, "id", id, "error", err)
			continue
		}

		files[id] = file
		files[id].Private = file.Private
	}

	return &api.ReadBatchFilesResponse{
		Files: model.MapFilesDictToProto(model.FileToProto, files),
	}, nil
}

func (h *Handler) ReadFile(ctx context.Context, req *api.ReadFileRequest) (*api.File, error) {
	logger.Debugw("read file", "user_id", req.UserId, "id", req.Id)

	file, err := h.repo.GetMeta(ctx, req.Id, "")
	if err != nil {
		return nil, err
	}

	return model.FileToProto(file), nil
}

type streamWriter struct {
	api.Files_ReadFileStreamServer
}

var _ io.Writer = (*streamWriter)(nil)

func (w *streamWriter) Write(p []byte) (n int, err error) {
	err = w.Send(&api.FileData{
		Data: &api.FileData_Chunk{
			Chunk: p,
		},
	})

	return len(p), err
}

func (h *Handler) ReadFileStream(params *api.ReadFileStreamRequest, stream api.Files_ReadFileStreamServer) (err error) {
	logger.Debugw("read file stream", "id", params.Id, "name", params.Name, "public", params.Public)

	file, err := h.repo.GetMeta(context.Background(), params.Id, params.Name)
	if err != nil {
		return err
	}

	if file.Private && params.Public {
		return errors.New("cannot read private file, when public requested")
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Id:             file.ID,
				UserId:         file.UserID,
				Name:           file.Name,
				Mime:           file.Mime,
				CreateUtcNano:  file.CreateUTCNano,
				Private:        file.Private,
				Size:           file.Size,
			},
		},
	})
	if err != nil {
		return
	}

	err = h.repo.ReadFile(context.Background(), file.OID, &streamWriter{stream})

	return
}

type streamReader struct {
	api.Files_SaveFileStreamServer
	mu  sync.Mutex
	buf bytes.Buffer
}

var _ io.Reader = (*streamReader)(nil)

func (r *streamReader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.buf.Len() > 0 {
		return r.buf.Read(p)
	}

	fileData, err := r.Recv()
	if err != nil {
		return 0, err
	}

	chunk, ok := fileData.Data.(*api.FileData_Chunk)
	if !ok {
		return 0, errors.New("file data chunk expected, wrong format")
	}

	_, err = r.buf.Write(chunk.Chunk)
	if err != nil {
		return 0, err
	}

	return r.buf.Read(p)
}

func (h *Handler) SaveFileStream(stream api.Files_SaveFileStreamServer) error {
	meta, err := stream.Recv()
	if err != nil {
		return err
	}

	file, ok := meta.Data.(*api.FileData_File)
	if !ok {
		return errors.New("wrong format: file meta expected")
	}

	logger.Debugw("save file stream", "id", file.File.Id, "user_id", file.File.UserId, "private", file.File.Private, "name", file.File.Name, "mime", file.File.Mime)

	err = h.repo.SaveFile(context.Background(), &streamReader{Files_SaveFileStreamServer: stream}, file.File.Id, file.File.UserId, file.File.Private, file.File.Name, file.File.Mime)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&api.SaveFileStreamResponse{})
}