package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/internal/domain"
)

func (m *DistributedMessages) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	timeout := 10*time.Second
	/* fsm.Apply() */
	future := m.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res = future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}

	return
}

func (m *DistributedMessages) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, userID int64, private bool, name string) (err error) {
	logger.Debugw("save message", "id", id, "text", text, "title", title, "file_ids", fileIDs, "user_id", userID, "private", private, "name", name)

	event, err := domain.CreateMessage(id, text, title, fileIDs, userID, private, name)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&AppendCommand{
		Id:       id,
		Text:     text,
		Title:    title,
		FileIds:  fileIDs,
		UserId:   userID,
		Private:  private,
		Name:     name,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)
	if err != nil {
		return
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *DistributedMessages) UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, userID int64) (err error) {
	logger.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "user_id", userID)

	event, err := domain.UpdateMessage(id, text, title, fileIDs, userID, name)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&UpdateCommand{
		Id:       id,
		UserId:   userID,
		FileIds:  fileIDs,
		Text:     text,
		Name:     name,
		Title:    title,
	})
	if err != nil {
		return err
	}

	res, err := m.apply(ctx, UpdateRequest, cmd)
	if err != nil {
		return err
	}

	switch updatedAt := res.(type) {
	case error:
		return updatedAt
	default:
		return m.publisher.Publish(context.Background(), event)
	}
}

func (m *DistributedMessages) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
	logger.Debugw("delete user messages", "user_id", userID)

	cmd, err := proto.Marshal(&DeleteUserMessagesCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteUserMessagesRequest, cmd)

	return
}

func (m *DistributedMessages) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("delete file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteFileCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteFileRequest, cmd)

	return
}

func (m *DistributedMessages) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	logger.Debugw("delete messages", "ids", ids, "user_id", userID)

	for _, id := range ids {

		event, err := domain.DeleteMessage(id, userID)
		if err != nil {
			return err
		}

		cmd, err := proto.Marshal(&DeleteCommand{
			Id:     id,
			UserId: userID,
		})
		if err != nil {
			return err
		}

		_, err = m.apply(ctx, DeleteRequest, cmd)
		if err != nil {
			return err
		}

		m.publisher.Publish(context.Background(), event)

	}

	return
}

func (m *DistributedMessages) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	logger.Debugw("publish messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&PublishCommand{
		Ids:           ids,
		UserId:        userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.PublishMessages(userID, ids)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *DistributedMessages) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	logger.Debugw("private messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateCommand{
		Ids:           ids,
		UserId:        userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.PrivateMessages(userID, ids)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

// TODO: pass one userID only, for public messages create ReadPublicMessage request
func (m *DistributedMessages) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *model.Message, err error) {
	logger.Debugw("read message", "id", id, "name", name, "user_ids", userIDs)

	message, err = m.messagesRepo.Read(ctx, userIDs, id, name)
	if err != nil {
		return
	}

	message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, id, userIDs)
	if err != nil {
		return
	}

	return
}

func (m *DistributedMessages) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
	logger.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)

	messages, isLastPage, err = m.messagesRepo.ReadMessages(ctx, userID, limit, offset)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, message.ID, []int64{userID})
		if err != nil {
			return
		}
	}

	return
}

func (m *DistributedMessages) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
	logger.Debugw("read batch messages", "user_id", userID, "ids", ids)

	messages, err = m.messagesRepo.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, message.ID, []int64{userID})
		if err != nil {
			return
		}
	}

	return
}
