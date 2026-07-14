package grpc

import (
	"io"
	"context"
	"bytes"
	"sync"
	"errors"

	"github.com/bd878/gallery/server/api/files"
)

type FilesController interface {
	ReadBatchFiles(ctx context.Context, userID int64, ids []int64) (dict map[int64]*files.File, err error)
	ReadFile(ctx context.Context, id int64, name string, public bool) (file *files.File, err error)
	ReadFileStream(ctx context.Context, oid int32, writer io.Writer) (err error)
	WriteFileStream(ctx context.Context, userID, id int64, private bool, name, description, mime string, reader io.Reader) (size int64, err error)
	ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*files.File, isLastPage bool, err error)
	PublishFiles(ctx context.Context, userID int64, ids []int64) (err error)
	PrivateFiles(ctx context.Context, userID int64, ids []int64) (err error)
	DeleteFiles(ctx context.Context, userID int64, ids []int64) (err error)
	ReadMessageFiles(ctx context.Context, id int64, userIDs []int64) (list []*files.File, err error)
}

type Handler struct {
	files.UnimplementedFilesServer
	controller FilesController
}

func NewFilesHandler(ctrl FilesController) *Handler {
	return &Handler{
		controller:  ctrl,
	}
}

func (h *Handler) ReadBatchFiles(ctx context.Context, req *files.ReadBatchFilesRequest) (resp *files.ReadBatchFilesResponse, err error) {
	dict, err := h.controller.ReadBatchFiles(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &files.ReadBatchFilesResponse{
		Files: dict,
	}

	return
}

func (h *Handler) ReadFile(ctx context.Context, req *files.ReadFileRequest) (*files.File, error) {
	file, err := h.controller.ReadFile(ctx, req.Id, req.Name, req.Public)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type streamWriter struct {
	files.Files_ReadFileStreamServer
}

var _ io.Writer = (*streamWriter)(nil)

func (w *streamWriter) Write(p []byte) (n int, err error) {
	err = w.Send(&files.FileData{
		Data: &files.FileData_Chunk{
			Chunk: p,
		},
	})

	return len(p), err
}

func (h *Handler) ReadFileStream(params *files.ReadFileStreamRequest, stream files.Files_ReadFileStreamServer) (err error) {
	file, err := h.controller.ReadFile(context.TODO(), params.Id, params.Name, params.Public)
	if err != nil {
		return err
	}

	err = stream.Send(&files.FileData{
		Data: &files.FileData_File{
			File: &files.File{
				Id:             file.Id,
				Oid:            file.Oid,
				UserId:         file.UserId,
				Name:           file.Name,
				Mime:           file.Mime,
				CreatedAt:      file.CreatedAt,
				UpdatedAt:      file.UpdatedAt,
				Private:        file.Private,
				Size:           file.Size,
			},
		},
	})
	if err != nil {
		return
	}

	return h.controller.ReadFileStream(context.TODO(), file.Oid, &streamWriter{stream})
}

type streamReader struct {
	files.Files_SaveFileStreamServer
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

	chunk, ok := fileData.Data.(*files.FileData_Chunk)
	if !ok {
		return 0, errors.New("file data chunk expected, wrong format")
	}

	_, err = r.buf.Write(chunk.Chunk)
	if err != nil {
		return 0, err
	}

	return r.buf.Read(p)
}

func (h *Handler) SaveFileStream(stream files.Files_SaveFileStreamServer) (err error) {
	meta, err := stream.Recv()
	if err != nil {
		return err
	}

	file, ok := meta.Data.(*files.FileData_File)
	if !ok {
		return errors.New("wrong format: file meta expected")
	}

	_, err = h.controller.WriteFileStream(context.TODO(), file.File.UserId, file.File.Id, file.File.Private,
		file.File.Name, file.File.Description, file.File.Mime, &streamReader{Files_SaveFileStreamServer: stream})
	if err != nil {
		return
	}

	return stream.SendAndClose(&files.SaveFileStreamResponse{})
}

func (h *Handler) ListFiles(ctx context.Context, req *files.ListFilesRequest) (resp *files.ListFilesResponse, err error) {
	list, isLastPage, err := h.controller.ListFiles(ctx, req.UserId, req.Limit, req.Offset, req.Asc, req.Private)
	if err != nil {
		return nil, err
	}

	resp = &files.ListFilesResponse{
		Files:      list,
		IsLastPage: isLastPage,
	}

	return
}

func (h *Handler) PublishFiles(ctx context.Context, req *files.PublishFilesRequest) (resp *files.PublishFilesResponse, err error) {
	err = h.controller.PublishFiles(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &files.PublishFilesResponse{}

	return
}

func (h *Handler) PrivateFiles(ctx context.Context, req *files.PrivateFilesRequest) (resp *files.PrivateFilesResponse, err error) {
	err = h.controller.PrivateFiles(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &files.PrivateFilesResponse{}

	return
}

func (h *Handler) DeleteFiles(ctx context.Context, req *files.DeleteFilesRequest) (resp *files.DeleteFilesResponse, err error) {
	err = h.controller.DeleteFiles(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &files.DeleteFilesResponse{}

	return
}

func (h *Handler) ReadMessageFiles(ctx context.Context, req *files.ReadMessageFilesRequest) (resp *files.ReadMessageFilesResponse, err error) {
	list, err := h.controller.ReadMessageFiles(ctx, req.Id, req.UserIds)
	if err != nil {
		return nil, err
	}

	resp = &files.ReadMessageFilesResponse{
		Files: list,
	}

	return
}
