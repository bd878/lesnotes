package application

import (
	"time"
	"context"
	"bytes"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/db/threads/pkg/machine"
	"github.com/bd878/gallery/server/db/threads/internal/controller"
)

type ThreadsRepository interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*threads.Thread, isLastPage bool, err error)
	ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool, privateMessage *bool) (ids []*threads.Thread, isLastPage bool, err error)
	ReadThreadByID(ctx context.Context, id, userID int64) (thread *threads.Thread, err error)
	ReadThreadByName(ctx context.Context, name string, userID int64) (thread *threads.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (path []*threads.PathStep, err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus       Consensus
	log             *logger.Logger
	threadsRepo     ThreadsRepository
}

func New(consensus Consensus, threadsRepo ThreadsRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:          log,
		consensus:    consensus,
		threadsRepo:  threadsRepo,
	}
}

func (m *Distributed) Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), duration)
}

func (m *Distributed) ResolveThread(ctx context.Context, id, userID int64) (path []*threads.PathStep, err error) {
	m.log.Debugw("resolve thread", "id", id, "user_id", userID)
	return m.threadsRepo.ResolveThread(ctx, id, userID)
}

func (m *Distributed) ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error) {
	m.log.Debugw("read thread", "id", id, "user_id", userID, "name", name)
	if name != "" {
		return m.threadsRepo.ReadThreadByName(ctx, name, userID)
	}
	return m.threadsRepo.ReadThreadByID(ctx, id, userID)
}

func (m *Distributed) ReadParent(ctx context.Context, id, userID int64) (parent *threads.Thread, err error) {
	m.log.Debugw("read parent", "id", id, "user_id", userID)

	thread, err := m.threadsRepo.ReadThreadByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if thread.ParentId == 0 {
		return nil, controller.ErrParentIsRoot
	}

	return m.threadsRepo.ReadThreadByID(ctx, thread.ParentId, userID)
}

func (m *Distributed) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*threads.Thread, isLastPage bool, err error) {
	m.log.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc)
	return m.threadsRepo.ListThreads(ctx, userID, parentID, limit, offset, asc)
}

func (m *Distributed) ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool, private *bool) (list []*threads.Thread, isLastPage bool, err error) {
	m.log.Debugw("list messages", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", "asc", asc, "private_message", private)
	return m.threadsRepo.ListMessages(ctx, userID, parentID, limit, offset, asc, private)
}

func (m *Distributed) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	m.log.Debugw("count threads", "user_id", userID, "id", id)
	return m.threadsRepo.CountThreads(ctx, id, userID)
}

func (m *Distributed) CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error) {
	m.log.Debugw("count messages", "user_id", userID, "id", id, "private_message", privateMessage)
	return m.threadsRepo.CountMessages(ctx, id, userID, privateMessage)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
