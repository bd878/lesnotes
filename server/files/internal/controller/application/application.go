package application

import (
	"io"
	"context"
	"errors"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/files/internal/domain"
)

type FilesRepository interface {
	SaveFile(ctx context.Context, reader io.Reader, userID, id int64, private bool, name, description, mime, createdAt, updatedAt string) (size int64, err error)
	GetMetaByID(ctx context.Context, id int64) (file *api.File, err error)
	GetMetaByName(ctx context.Context, fileName string) (file *api.File, err error)
	DeleteFiles(ctx context.Context, userID int64, ids []int64) (err error)
	ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error)
	ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*api.File, isLastPage bool, err error)
	PublishFiles(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	PrivateFiles(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
}

type Application struct {
	log              *logger.Logger
	publisher        ddd.EventPublisher[ddd.Event]
	filesRepo        FilesRepository
}

func New(publisher ddd.EventPublisher[ddd.Event], filesRepo FilesRepository, log *logger.Logger) *Application {
	return &Application{
		log:        log,
		publisher:  publisher,
		filesRepo:  filesRepo,
	}
}

func (a *Application) ReadBatchFiles(ctx context.Context, userID int64, ids []int64) (files map[int64]*api.File, err error) {
	a.log.Debugw("read batch files", "user_id", userID, "ids", ids)

	files = make(map[int64]*api.File, len(ids))
	for _, id := range ids {
		file, err := a.filesRepo.GetMetaByID(ctx, id)
		if err != nil {
			files[id] = &api.File{Error: err.Error()}
			logger.Errorw("failed to read file", "user_id", userID, "id", id, "error", err)
			continue
		}

		files[id] = file
	}

	return files, nil
}

func (a *Application) ReadFile(ctx context.Context, id int64, name string, public bool) (file *api.File, err error) {
	a.log.Debugw("read file", "id", id, "name", name, "public", public)

	if name != "" {
		file, err = a.filesRepo.GetMetaByName(ctx, name)
	} else {
		file, err = a.filesRepo.GetMetaByID(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	if file.Private && public {
		return nil, errors.New("cannot read private file")
	}

	return file, nil
}

func (a *Application) ReadFileStream(ctx context.Context, oid int32, writer io.Writer) (err error) {
	a.log.Debugw("read file stream", "oid", oid)

	return a.filesRepo.ReadFile(ctx, oid, writer)
}

func (a *Application) WriteFileStream(ctx context.Context, userID, id int64, private bool, name, description, mime string, reader io.Reader) (size int64, err error) {
	a.log.Debugw("write file stream", "user_id", userID, "id", id, "private", private, "name", name, "description", description, "mime", mime)

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	size, err = a.filesRepo.SaveFile(context.TODO(), reader, userID, id, private, name, description, mime, createdAt, updatedAt)
	if err != nil {
		return
	}

	event, err := domain.UploadFile(id, name, description, userID, private, mime, size, createdAt, updatedAt)
	if err != nil {
		return 0, err
	}

	err = a.publisher.Publish(context.TODO(), event)
	if err != nil {
		return
	}

	return
}

func (a *Application) ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*api.File, isLastPage bool, err error) {
	a.log.Debugw("list files", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending, "private", private)

	return a.filesRepo.ListFiles(ctx, userID, limit, offset, ascending, private)
}

func (a *Application) PublishFiles(ctx context.Context, userID int64, ids []int64) (err error) {
	a.log.Debugw("publish files", "user_id", userID, "ids", ids)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.PublishFiles(userID, ids, updatedAt)
	if err != nil {
		return err
	}

	err = a.filesRepo.PublishFiles(ctx, userID, ids, updatedAt)
	if err != nil {
		return
	}

	return a.publisher.Publish(context.TODO(), event)
}

func (a *Application) PrivateFiles(ctx context.Context, userID int64, ids []int64) (err error) {
	a.log.Debugw("private files", "user_id", userID, "ids", ids)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.PrivateFiles(userID, ids, updatedAt)
	if err != nil {
		return err
	}

	err = a.filesRepo.PrivateFiles(ctx, userID, ids, updatedAt)
	if err != nil {
		return
	}

	return a.publisher.Publish(context.TODO(), event)
}

func (a *Application) DeleteFiles(ctx context.Context, userID int64, ids []int64) (err error) {
	a.log.Debugw("delete files", "user_id", userID, "ids", ids)

	event, err := domain.DeleteFiles(userID, ids)
	if err != nil {
		return err
	}

	err = a.filesRepo.DeleteFiles(ctx, userID, ids)
	if err != nil {
		return
	}

	return a.publisher.Publish(context.TODO(), event)
}