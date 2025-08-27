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

func (m *DistributedMessages) SaveMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (err error) {
	logger.Debugw("save message", "id", id, "text", text, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private, "name", name)

	cmd, err := proto.Marshal(&AppendCommand{
		Id:       id,
		Text:     text,
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

func (m *DistributedMessages) UpdateMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private int32) (err error) {
	cmd, err := proto.Marshal(&UpdateCommand{
		Id:       id,
		UserId:   userID,
		FileIds:  fileIDs,
		ThreadId: threadID,
		Text:     text,
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
	return m.repo.Read(ctx, userIDs, id)
}

func (m *DistributedMessages) ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool, private int32) (messages []*model.Message, isLastPage bool, err error) {
	return m.repo.ReadThreadMessages(ctx, userID, threadID, limit, offset, private)
}

func (m *DistributedMessages) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool, private int32) (messages []*model.Message, isLastPage bool, err error) {
	return m.repo.ReadMessages(ctx, userID, limit, offset, private)
}
