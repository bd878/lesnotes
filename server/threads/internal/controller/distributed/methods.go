package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/threads/internal/domain"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

func (m *DistributedThreads) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
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

func (m *DistributedThreads) CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error) {
	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID,
		"next_id", nextID, "prev_id", prevID, "name", name, "description", description, "private", private)

	event, err := domain.CreateThread(id, userID, parentID, name, description, private)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&AppendCommand{
		Id:       id,
		UserId:   userID,
		ParentId: parentID,
		NextId:   nextID,
		PrevId:   prevID,
		Name:     name,
		Private:  private,
		Description: description,
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

func (m *DistributedThreads) UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error) {
	logger.Debugw("update thread", "id", id, "user_id", userID, "name", name, "description", description)

	event, err := domain.UpdateThread(id, userID, name, description)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&UpdateCommand{
		Id:           id,
		UserId:       userID,
		Name:         name,
		Description:  description,
	})
	if err != nil {
		return err
	}

	res, err := m.apply(ctx, UpdateRequest, cmd)
	if err != nil {
		return err
	}

	switch val := res.(type) {
	case error:
		return val
	default:
		return m.publisher.Publish(context.Background(), event)
	}
}

func (m *DistributedThreads) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	logger.Debugw("reorder thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID)

	cmd, err := proto.Marshal(&ReorderCommand{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		NextId:    nextID,
		PrevId:    prevID,
	})
	if err != nil {
		logger.Debugw("failed to marshal", "error", err)
		return err
	}

	_, err = m.apply(ctx, ReorderRequest, cmd)
	if err != nil {
		return
	}

	if parentID != -1 {
		event, err := domain.ChangeThreadParent(id, userID, parentID)
		if err != nil {
			return err
		}

		err = m.publisher.Publish(context.Background(), event)
	}

	return
}

func (m *DistributedThreads) PrivateThread(ctx context.Context, id, userID int64) error {
	logger.Debugw("private thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.PrivateThread(id, userID)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *DistributedThreads) PublishThread(ctx context.Context, id int64, userID int64) error {
	logger.Debugw("publich thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PublishCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.PublishThread(id, userID)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *DistributedThreads) DeleteThread(ctx context.Context, id, userID int64) error {
	logger.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.DeleteThread(id, userID)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *DistributedThreads) ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error) {
	logger.Debugw("resolve thread", "id", id, "user_id", userID)
	return m.repo.ResolveThread(ctx, id, userID)
}

func (m *DistributedThreads) ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error) {
	logger.Debugw("read thread", "id", id, "user_id", userID, "name", name)
	return m.repo.ReadThread(ctx, id, userID, name)
}

func (m *DistributedThreads) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*threads.Thread, isLastPage bool, err error) {
	logger.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc)
	return m.repo.ListThreads(ctx, userID, parentID, limit, offset, asc)
}

func (m *DistributedThreads) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	logger.Debugw("count threads", "user_id", userID, "id", id)
	return m.repo.CountThreads(ctx, id, userID)
}
