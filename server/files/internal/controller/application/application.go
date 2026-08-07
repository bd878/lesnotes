package application

import (
	"context"
	"fmt"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/bd878/gallery/server/api/files"
	"github.com/bd878/gallery/server/files/internal/domain"
	"github.com/bd878/gallery/server/internal/ddd"
)

type FilesRepository interface {
	SaveFile(ctx context.Context, reader io.Reader, userID, id int64, private bool, name, description, mime, createdAt, updatedAt string) (size int64, err error)
	GetMetaByID(ctx context.Context, id int64) (file *files.File, err error)
	GetMetaByName(ctx context.Context, fileName string) (file *files.File, err error)
	DeleteFiles(ctx context.Context, userID int64, ids []int64) (err error)
	ReadFile(ctx context.Context, oid int32, writer io.Writer) (err error)
	ReadBatchFiles(ctx context.Context, ids []int64) (list []*files.File, err error)
	ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*files.File, isLastPage bool, err error)
	PublishFiles(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	PrivateFiles(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
}

type MessagesRepository interface {
	SaveMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) (err error)
	UpdateMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) (err error)
	ReadMessageFiles(ctx context.Context, id int64, userIDs []int64) (fileIDs []int64, err error)
	DeleteFiles(ctx context.Context, ids []int64) (err error)
	DeleteMessage(ctx context.Context, id, userID int64) (err error)
}

type Application struct {
	publisher    ddd.EventPublisher[ddd.Event]
	filesRepo    FilesRepository
	messagesRepo MessagesRepository
}

func New(publisher ddd.EventPublisher[ddd.Event],
	filesRepo FilesRepository, messagesRepo MessagesRepository) *Application {
	return &Application{
		publisher:    publisher,
		filesRepo:    filesRepo,
		messagesRepo: messagesRepo,
	}
}

func (a *Application) SaveMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) (err error) {
	slog.Debug("save message files", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("file_ids", fmt.Sprintf("%v", fileIDs)))

	return a.messagesRepo.SaveMessageFiles(ctx, id, userID, fileIDs)
}

func (a *Application) DeleteMessageFiles(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("delete message files", slog.Int64("id", id), slog.Int64("user_id", userID))

	fileIDs, err := a.messagesRepo.ReadMessageFiles(ctx, id, []int64{userID})
	if err != nil {
		slog.Debug("no files")
		return err
	}

	if len(fileIDs) == 0 {
		return nil
	}

	return a.DeleteFiles(ctx, userID, fileIDs)
}

func (a *Application) UpdateMessageFiles(ctx context.Context, id, userID int64, fileIDs []int64) (err error) {
	slog.Debug("update message files", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("file_ids", fmt.Sprintf("%v", fileIDs)))
	return a.messagesRepo.UpdateMessageFiles(ctx, id, userID, fileIDs)
}

func (a *Application) ReadMessageFiles(ctx context.Context, id int64, userIDs []int64) (list []*files.File, err error) {
	slog.Debug("read message files", slog.Int64("id", id), slog.String("user_ids", fmt.Sprintf("%v", userIDs)))

	fileIDs, err := a.messagesRepo.ReadMessageFiles(ctx, id, userIDs)
	if err != nil {
		return nil, err
	}

	return a.filesRepo.ReadBatchFiles(ctx, fileIDs)
}

func (a *Application) PublishMessageFiles(ctx context.Context, userID int64, messageIDs []int64) (err error) {
	slog.Debug("publish message files", slog.Int64("user_id", userID), slog.String("message_ids", fmt.Sprintf("%v", messageIDs)))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	for _, messageID := range messageIDs {
		fileIDs, err := a.messagesRepo.ReadMessageFiles(ctx, messageID, []int64{userID})
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		err = a.filesRepo.PublishFiles(ctx, userID, fileIDs, updatedAt)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}

	return nil
}

func (a *Application) PrivateMessageFiles(ctx context.Context, userID int64, messageIDs []int64) (err error) {
	slog.Debug("private message files", slog.Int64("user_id", userID), slog.String("message_ids", fmt.Sprintf("%v", messageIDs)))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	for _, messageID := range messageIDs {
		fileIDs, err := a.messagesRepo.ReadMessageFiles(ctx, messageID, []int64{userID})
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		err = a.filesRepo.PrivateFiles(ctx, userID, fileIDs, updatedAt)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}

	return nil
}

func (a *Application) ReadBatchFiles(ctx context.Context, userID int64, ids []int64) (list map[int64]*files.File, err error) {
	slog.Debug("read batch files", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

	list = make(map[int64]*files.File, len(ids))
	for _, id := range ids {
		file, err := a.filesRepo.GetMetaByID(ctx, id)
		if err != nil {
			list[id] = &files.File{Error: err.Error()}
			slog.Error(err.Error())
			continue
		}

		list[id] = file
	}

	return list, nil
}

func (a *Application) ReadFile(ctx context.Context, id int64, name string, public bool) (file *files.File, err error) {
	slog.Debug("read file", slog.Int64("id", id), slog.String("name", name), slog.Bool("public", public))

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
	slog.Debug("read file stream", slog.Int("oid", int(oid)))
	return a.filesRepo.ReadFile(ctx, oid, writer)
}

func (a *Application) WriteFileStream(ctx context.Context, userID, id int64, private bool, name, description, mime string, reader io.Reader) (size int64, err error) {
	slog.Debug("write file stream",
		slog.Int64("user_id", userID),
		slog.Int64("id", id),
		slog.Bool("private", private),
		slog.String("name", name),
		slog.String("description", description),
		slog.String("mime", mime),
	)

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

func (a *Application) ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list []*files.File, isLastPage bool, err error) {
	slog.Debug("list files", slog.Int64("user_id", userID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("private", private))
	return a.filesRepo.ListFiles(ctx, userID, limit, offset, ascending, private)
}

func (a *Application) PublishFiles(ctx context.Context, userID int64, ids []int64) (err error) {
	slog.Debug("publish files", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

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
	slog.Debug("private files", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

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
	slog.Debug("delete files", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

	event, err := domain.DeleteFiles(userID, ids)
	if err != nil {
		return err
	}

	err = a.filesRepo.DeleteFiles(ctx, userID, ids)
	if err != nil {
		return
	}

	err = a.messagesRepo.DeleteFiles(ctx, ids)
	if err != nil {
		return
	}

	return a.publisher.Publish(context.TODO(), event)
}
