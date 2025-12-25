package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type RepoConnection interface {
	Release()
}

type Repository interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*threads.Thread, isLastPage bool, err error)
	ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error)
	AppendThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error)
	UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error)
	PrivateThread(ctx context.Context, id, userID int64) error
	PublishThread(ctx context.Context, id, userID int64) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	Truncate(ctx context.Context) error
	Dump(ctx context.Context) (reader io.ReadCloser, err error)
	Restore(ctx context.Context, reader io.ReadCloser) (err error)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	repo Repository
	// publisher       ddd.EventPublisher[ddd.Event]
}

func (f *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendRequest:
		return f.applyAppend(buf[1:])
	case UpdateRequest:
		return f.applyUpdate(buf[1:])
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case PublishRequest:
		return f.applyPublish(buf[1:])
	case PrivateRequest:
		return f.applyPrivate(buf[1:])
	case ReorderRequest:
		return f.applyReorder(buf[1:])
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.AppendThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId, cmd.PrevId, cmd.Name, cmd.Description, cmd.Private)
	// TODO: left to distribute on prod existing events
	// if err != nil {
	// 	return err
	// }

	// event, err := domain.CreateThread(cmd.Id, cmd.UserId, cmd.ParentId, cmd.Name, cmd.Description, cmd.Private)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugln("publish create thread")

	// return f.publisher.Publish(context.Background(), event)
}

func (f *fsm) applyReorder(raw []byte) interface{} {
	var cmd ReorderCommand
	proto.Unmarshal(raw, &cmd)

	logger.Debugln("apply reorder")

	return f.repo.ReorderThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId, cmd.PrevId)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugw("repo reorder thread error", "error", err)

	// if cmd.ParentId != -1 {
	// 	logger.Debugln("publish reorder thread")

	// 	event, err := domain.ChangeThreadParent(cmd.Id, cmd.UserId, cmd.ParentId)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return f.publisher.Publish(context.Background(), event)
	// }

	// return nil
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.UpdateThread(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugln("publish update thread")

	// event, err := domain.UpdateThread(cmd.Id, cmd.UserId, cmd.Name, cmd.Description)
	// if err != nil {
	// 	return err
	// }

	// return f.publisher.Publish(context.Background(), event)
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.DeleteThread(context.Background(), cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugln("publish delete thread")

	// event, err := domain.DeleteThread(cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// return f.publisher.Publish(context.Background(), event)
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.PublishThread(context.Background(), cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugln("publish publish thread")

	// event, err := domain.PublishThread(cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// return f.publisher.Publish(context.Background(), event)
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.PrivateThread(context.Background(), cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// logger.Debugln("publish private thread")

	// event, err := domain.PrivateThread(cmd.Id, cmd.UserId)
	// if err != nil {
	// 	return err
	// }

	// return f.publisher.Publish(context.Background(), event)
}
