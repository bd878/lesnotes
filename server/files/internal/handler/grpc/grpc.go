package grpc

import (
	"io"
	"context"
	"bytes"
	"sync"
	"errors"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/files/pkg/model"
	"github.com/bd878/gallery/server/files/internal/domain"
)

// TODO: add controller, derive domain events on controller level

type Repository interface {
	SaveFile(ctx context.Context, reader io.Reader, id, userID int64, private bool, name, description, mime string) (size int64, err error)
	GetMeta(ctx context.Context, id int64, fileName string) (file *model.File, err error)
	DeleteFile(ctx context.Context, id, ownerID int64) (err error)
	ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error)
	ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*model.File, isLastPage bool, err error)
	PublishFile(ctx context.Context, id, userID int64) (err error)
	PrivateFile(ctx context.Context, id, userID int64) (err error)
}

type Handler struct {
	api.UnimplementedFilesServer
	repo       Repository
	publisher  ddd.EventPublisher[ddd.Event]
}

func New(repo Repository, publisher ddd.EventPublisher[ddd.Event]) *Handler {
	return &Handler{
		repo:      repo,
		publisher: publisher,
	}
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
	}

	return &api.ReadBatchFilesResponse{
		Files: model.MapFilesDictToProto(model.FileToProto, files),
	}, nil
}

func (h *Handler) ReadFile(ctx context.Context, req *api.ReadFileRequest) (*api.File, error) {
	logger.Debugw("read file", "user_id", req.UserId, "id", req.Id, "public", req.Public)

	file, err := h.repo.GetMeta(ctx, req.Id, "")
	if err != nil {
		return nil, err
	}

	if file.Private && file.UserID != req.UserId {
		return nil, errors.New("requested file is private")
	}

	if file.Private && req.Public {
		return nil, errors.New("cannot read private file")
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

func (h *Handler) SaveFileStream(stream api.Files_SaveFileStreamServer) (err error) {
	meta, err := stream.Recv()
	if err != nil {
		return err
	}

	file, ok := meta.Data.(*api.FileData_File)
	if !ok {
		return errors.New("wrong format: file meta expected")
	}

	logger.Debugw("save file stream", "id", file.File.Id, "user_id", file.File.UserId, "private", file.File.Private, "name", file.File.Name, "description", file.File.Description, "mime", file.File.Mime)

	var size int64
	size, err = h.repo.SaveFile(context.Background(), &streamReader{Files_SaveFileStreamServer: stream}, file.File.Id, file.File.UserId, file.File.Private, file.File.Name, file.File.Description, file.File.Mime)
	if err != nil {
		return err
	}

	event, err := domain.UploadFile(file.File.Id, file.File.Name, file.File.Description, file.File.UserId, file.File.Private, file.File.Mime, size)
	if err != nil {
		return err
	}

	err = stream.SendAndClose(&api.SaveFileStreamResponse{})
	if err != nil {
		return
	}

	return h.publisher.Publish(context.Background(), event)
}

func (h *Handler) ListFiles(ctx context.Context, req *api.ListFilesRequest) (resp *api.ListFilesResponse, err error) {
	logger.Debugw("list files", "user_id", req.UserId, "limit", req.Limit, "offset", req.Offset, "asc", req.Asc, "private", req.Private)

	list, isLastPage, err := h.repo.ListFiles(context.Background(), req.UserId, req.Limit, req.Offset, req.Asc, req.Private)
	if err != nil {
		return nil, err
	}

	resp = &api.ListFilesResponse{
		Files:      model.MapFilesToProto(model.FileToProto, list),
		IsLastPage: isLastPage,
	}

	return
}

func (h *Handler) PublishFile(ctx context.Context, req *api.PublishFileRequest) (resp *api.PublishFileResponse, err error) {
	logger.Debugw("publish file", "id", req.Id, "user_id", req.UserId)

	event, err := domain.PublishFile(req.UserId, req.Id)
	if err != nil {
		return nil, err
	}

	err = h.repo.PublishFile(ctx, req.Id, req.UserId)
	if err != nil {
		return
	}

	err = h.publisher.Publish(context.Background(), event)
	if err != nil {
		return
	}

	resp = &api.PublishFileResponse{}

	return
}

func (h *Handler) PrivateFile(ctx context.Context, req *api.PrivateFileRequest) (resp *api.PrivateFileResponse, err error) {
	logger.Debugw("private file", "id", req.Id, "user_id", req.UserId)

	event, err := domain.PrivateFile(req.UserId, req.Id)
	if err != nil {
		return nil, err
	}

	err = h.repo.PrivateFile(ctx, req.Id, req.UserId)
	if err != nil {
		return
	}

	err = h.publisher.Publish(context.Background(), event)
	if err != nil {
		return
	}

	resp = &api.PrivateFileResponse{}

	return
}

func (h *Handler) DeleteFile(ctx context.Context, req *api.DeleteFileRequest) (resp *api.DeleteFileResponse, err error) {
	logger.Debugw("delete file", "id", req.Id, "user_id", req.UserId)

	event, err := domain.DeleteFile(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	err = h.repo.DeleteFile(ctx, req.UserId, req.Id)
	if err != nil {
		return
	}

	err = h.publisher.Publish(context.Background(), event)
	if err != nil {
		return
	}

	resp = &api.DeleteFileResponse{}

	return
}