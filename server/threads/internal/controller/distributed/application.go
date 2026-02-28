package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/threads/internal/domain"
	"github.com/bd878/gallery/server/threads/internal/machine"
	"github.com/bd878/gallery/server/threads/pkg/model"
)

type ThreadsRepository interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*model.Thread, isLastPage bool, err error)
	ReadThreadByID(ctx context.Context, id, userID int64) (thread *model.Thread, err error)
	ReadThreadByName(ctx context.Context, name string, userID int64) (thread *model.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus       Consensus
	log             *logger.Logger
	threadsRepo     ThreadsRepository
	publisher       ddd.EventPublisher[ddd.Event]
}

func New(consensus Consensus, publisher ddd.EventPublisher[ddd.Event], threadsRepo ThreadsRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:          log,
		publisher:    publisher,
		consensus:    consensus,
		threadsRepo:  threadsRepo,
	}
}

func (m *Distributed) apply(ctx context.Context, reqType machine.RequestType, cmd []byte) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), 10*time.Second)
}

func (m *Distributed) CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error) {
	m.log.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID,
		"next_id", nextID, "prev_id", prevID, "name", name, "description", description, "private", private)


	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.CreateThread(id, userID, parentID, name, description, private, createdAt, updatedAt)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.AppendCommand{
		Id:            id,
		UserId:        userID,
		ParentId:      parentID,
		NextId:        nextID,
		PrevId:        prevID,
		Name:          name,
		Private:       private,
		Description:   description,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendRequest, cmd)
	if err != nil {
		return
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) UpdateThread(ctx context.Context, id, userID int64, name, description *string) (err error) {
	m.log.Debugw("update thread", "id", id, "user_id", userID, "name", name, "description", description)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateThread(id, userID, name, description, updatedAt)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.UpdateCommand{
		Id:            id,
		UserId:        userID,
		Name:          name,
		Description:   description,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateRequest, cmd)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	m.log.Debugw("reorder thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID)

	cmd, err := proto.Marshal(&machine.ReorderCommand{
		Id:            id,
		UserId:        userID,
		ParentId:      parentID,
		NextId:        nextID,
		PrevId:        prevID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.ReorderRequest, cmd)
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

func (m *Distributed) PrivateThread(ctx context.Context, id, userID int64) error {
	m.log.Debugw("private thread", "id", id, "user_id", userID)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&machine.PrivateCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PrivateRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.PrivateThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) PublishThread(ctx context.Context, id int64, userID int64) error {
	m.log.Debugw("publich thread", "id", id, "user_id", userID)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&machine.PublishCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PublishRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.PublishThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) DeleteThread(ctx context.Context, id, userID int64) error {
	m.log.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteCommand{
		Id:       id,
		UserId:   userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteRequest, cmd)
	if err != nil {
		return err
	}

	event, err := domain.DeleteThread(id, userID)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error) {
	m.log.Debugw("resolve thread", "id", id, "user_id", userID)
	return m.threadsRepo.ResolveThread(ctx, id, userID)
}

func (m *Distributed) ReadThread(ctx context.Context, id, userID int64, name string) (thread *model.Thread, err error) {
	m.log.Debugw("read thread", "id", id, "user_id", userID, "name", name)
	if name != "" {
		return m.threadsRepo.ReadThreadByName(ctx, name, userID)
	}
	return m.threadsRepo.ReadThreadByID(ctx, id, userID)
}

func (m *Distributed) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*model.Thread, isLastPage bool, err error) {
	m.log.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc)
	return m.threadsRepo.ListThreads(ctx, userID, parentID, limit, offset, asc)
}

func (m *Distributed) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	m.log.Debugw("count threads", "user_id", userID, "id", id)
	return m.threadsRepo.CountThreads(ctx, id, userID)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
