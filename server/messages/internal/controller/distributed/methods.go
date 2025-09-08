package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
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

func (m *DistributedMessages) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (err error) {
	logger.Debugw("save message", "id", id, "text", text, "title", title, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private, "name", name)

	cmd, err := proto.Marshal(&AppendCommand{
		Id:       id,
		Text:     text,
		Title:    title,
		FileIds:  fileIDs,
		ThreadId: threadID,
		UserId:   userID,
		Private:  private,
		Name:     name,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)

	return
}

func (m *DistributedMessages) UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, threadID int64, userID int64, private int32) (err error) {
	logger.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private)

	cmd, err := proto.Marshal(&UpdateCommand{
		Id:       id,
		UserId:   userID,
		FileIds:  fileIDs,
		ThreadId: threadID,
		Text:     text,
		Name:     name,
		Title:    title,
		Private:  private,
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
		return nil
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

func (m *DistributedMessages) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	logger.Debugw("delete messages", "ids", ids, "user_id", userID)

	for _, id := range ids {
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

	return
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

	return
}

func (m *DistributedMessages) ReadMessage(ctx context.Context, id int64, userIDs []int64) (message *model.Message, err error) {
	logger.Debugw("read message", "id", id, "user_ids", userIDs)
	return m.repo.Read(ctx, userIDs, id)
}

func (m *DistributedMessages) ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
	logger.Debugw("read thread messages", "user_id", userID, "thread_id", threadID, "limit", limit, "offset", offset, "ascending", ascending)
	return m.repo.ReadThreadMessages(ctx, userID, threadID, limit, offset)
}

func (m *DistributedMessages) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
	logger.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)
	return m.repo.ReadMessages(ctx, userID, limit, offset)
}

func (m *DistributedMessages) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
	logger.Debugw("read batch messages", "user_id", userID, "ids", ids)
	return m.repo.ReadBatchMessages(ctx, userID, ids)
}

func (m *DistributedMessages) ReadPath(ctx context.Context, userID, id int64) (path []*model.Message, err error) {
	logger.Debugw("read path", "user_id", userID, "id", id)
	return m.repo.ReadPath(ctx, userID, id)
}