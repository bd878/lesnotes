package grpc

import (
	"io"
	"time"
	"context"
	"bytes"
	"sync"
	"errors"

	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Repository interface {
	SaveFile(ctx context.Context, reader io.Reader, file *model.File) (err error)
	GetMeta(ctx context.Context, ownerID, id int32) (file *model.File, err error)
	DeleteFile(ctx context.Context, ownerID, id int32) (err error)
	ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error)
}

type Handler struct {
	api.UnimplementedFilesServer
	repo       Repository
}

func New(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ReadBatchFiles(ctx context.Context, req *api.ReadBatchFilesRequest) (
	*api.ReadBatchFilesResponse, error,
) {
	files := make(map[int32]*model.File, len(req.Ids))
	for _, id := range req.Ids {
		files[id] = &model.File{
			ID:     id,
			UserID: req.UserId,
		}

		file, err := h.repo.GetMeta(ctx, req.UserId, id)
		if err != nil {
			files[id].Error = "can not found file"
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

func (h *Handler) ReadFile(ctx context.Context, req *api.ReadFileRequest) (
	*api.File, error,
) {
	file, err := h.repo.GetMeta(ctx, req.UserId, req.Id)
	if err != nil {
		logger.Errorw("failed to read one file", "user_id", req.UserId, "file_id", req.Id, "error", err)
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
	file, err := h.repo.GetMeta(context.Background(), params.UserId, params.Id)
	if err != nil {
		logger.Errorw("failed to read file", "user_id", params.UserId, "id", params.Id, "name", params.Name, "public", params.Public, "error", err)
		return err
	}

	if file.Private && params.Public {
		logger.Errorw("failed to read private file", "user_id", params.UserId, "id", params.Id, "name", params.Name, "public", params.Public)
		return errors.New("cannot read private file, when public requested")
	}

	err = stream.Send(&api.FileData{
		Data: &api.FileData_File{
			File: &api.File{
				Id:             file.ID,
				UserId:         file.UserID,
				Name:           file.Name,
				CreateUtcNano:  file.CreateUTCNano,
				Private:        file.Private,
				Size:           file.Size,
			},
		},
	})
	if err != nil {
		logger.Errorw("stream failed to send filedata", "user_id", file.UserID, "id", file.ID, "error", err)
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
		logger.Errorw("save file stream failed to receive meta", "error", err)
		return err
	}

	file, ok := meta.Data.(*api.FileData_File)
	if !ok {
		logger.Errorw("send file data first, then chunk", "error", "wrong format")
		return errors.New("wrong format: file meta expected")
	}

	id := utils.RandomID()
	timeCreated := time.Now().UnixNano()

	err = h.repo.SaveFile(context.Background(), &streamReader{Files_SaveFileStreamServer: stream}, &model.File{
		ID:              id,
		UserID:          file.File.UserId,
		Name:            file.File.Name,
		CreateUTCNano:   timeCreated,
		Private:         file.File.Private,
	})
	if err != nil {
		logger.Errorw("failed to save file meta", "user_id", file.File.UserId, "name", file.File.Name, "error", err)
		return err
	}

	return stream.SendAndClose(&api.SaveFileStreamResponse{
		File: &api.File{
			Id:               id,
			Name:             file.File.Name,
			CreateUtcNano:    timeCreated,
		},
	})
}