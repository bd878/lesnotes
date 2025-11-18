package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
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

func (m *Distributed) CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name string, private bool) (err error) {
	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID, "name", name, "private", private)

	cmd, err := proto.Marshal(&AppendCommand{
		Id:       id,
		UserId:   userID,
		ParentId: parentID,
		NextId:   nextID,
		PrevId:   prevID,
		Name:     name,
		Private:  private,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)

	return
}

func (m *Distributed) UpdateThread(ctx context.Context, id, userID int64, name string, private int32) (err error) {
	logger.Debugw("update thread", "id", id, "user_id", userID, "name", name, "private", private)

	cmd, err := proto.Marshal(&UpdateCommand{
		Id:       id,
		UserId:   userID,
		Name:     name,
		Private:  private,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, UpdateRequest, cmd)

	return
}

func (m *Distributed) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	logger.Debugw("reorder thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID)

	cmd, err := proto.Marshal(&ReorderCommand{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		NextId:    nextID,
		PrevId:    prevID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, ReorderRequest, cmd)

	return
}

func (m *Distributed) PrivateThread(ctx context.Context, id, userID int64) error {
	logger.Debugw("private thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateRequest, cmd)

	return err
}

func (m *Distributed) PublishThread(ctx context.Context, id int64, userID int64) error {
	logger.Debugw("publich thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PublishCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishRequest, cmd)

	return err
}

func (m *Distributed) DeleteThread(ctx context.Context, id, userID int64) error {
	logger.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteRequest, cmd)

	return err
}

func (m *Distributed) ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error) {
	logger.Debugw("resolve thread", "id", id, "user_id", userID)
	return m.repo.ResolveThread(ctx, id, userID)
}

func (m *Distributed) ReadThread(ctx context.Context, id, userID int64) (thread *threads.Thread, err error) {
	logger.Debugw("read thread", "id", id, "user_id", userID)
	return m.repo.ReadThread(ctx, id, userID)
}