package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/search/pkg/model"
)

func (m *Distributed) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
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

func (m *Distributed) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
	// for integration events; though raft will not allow .apply for not a leader, anyway
	if !m.isLeader() {
		return
	}

	logger.Debugw("save search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text, "private", private)

	cmd, err := proto.Marshal(&AppendCommand{
		Id:      id,
		Text:    text,
		Title:   title,
		Name:    name,
		UserId:  userID,
		Private: private,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)
	if err != nil {
		return
	}

	return
}

func (m *Distributed) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return
	}

	logger.Debugw("delete search message", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteRequest, cmd)

	return
}

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string) (list []*model.Message, err error) {
	logger.Debugw("search messages", "user_id", userID, "substr", substr)
	return m.repo.SearchMessages(ctx, userID, substr)
}