package service

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/db/threads/pkg/machine"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/threads/internal/domain"
	"github.com/bd878/gallery/server/threads/pkg/model"
)

type Controller struct {
	client    threads.ThreadsClient
	publisher ddd.EventPublisher[ddd.Event]
}

func New(container di.Container, publisher ddd.EventPublisher[ddd.Event]) *Controller {
	client := container.Get("threadsClient").(threads.ThreadsClient)

	return &Controller{
		client:    client,
		publisher: publisher,
	}
}

func (s *Controller) ReadThread(ctx context.Context, id, userID int64, name string) (thread *model.Thread, err error) {
	slog.Debug("read thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("name", name))

	resp, err := s.client.Read(ctx, &threads.ReadRequest{Id: id, UserId: userID, Name: name})
	if err != nil {
		return nil, err
	}

	thread = model.ThreadFromProto(resp)

	return
}

func (s *Controller) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*model.Thread, isLastPage bool, err error) {
	slog.Debug("list threads", slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("asc", asc))

	resp, err := s.client.List(ctx, &threads.ListRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
		Asc:      asc,
	})
	if err != nil {
		return nil, false, err
	}

	return model.MapThreadsFromProto(model.ThreadFromProto, resp.List), resp.IsLastPage, err
}

func (s *Controller) ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool, privateMessage *bool) (list []*model.Thread, isLastPage bool, err error) {
	slog.Debug("list messages", slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("asc", asc), slog.Any("private_message", privateMessage))

	resp, err := s.client.ListMessages(ctx, &threads.ListMessagesRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
		Asc:      asc,
		Private:  privateMessage,
	})
	if err != nil {
		return nil, false, err
	}

	return model.MapThreadsFromProto(model.ThreadFromProto, resp.List), resp.IsLastPage, err
}

func (s *Controller) ResolveThread(ctx context.Context, id, userID int64) (path []*model.PathStep, err error) {
	slog.Debug("resolve thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	resp, err := s.client.Resolve(ctx, &threads.ResolveRequest{Id: id, UserId: userID})
	if err != nil {
		return nil, err
	}

	path = model.PathFromProto(resp.Path)

	return
}

func (s *Controller) PublishThread(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("publish thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.PublishThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&threads.PublishCommand{
		Id:        id,
		UserId:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.PublishRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("private thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.PrivateThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&threads.PrivateCommand{
		Id:        id,
		UserId:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.PrivateRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return 
}

func (s *Controller) CreateThread(ctx context.Context, id, userID, parentID int64, name, description, title string, private bool) (err error) {
	slog.Debug("create thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.String("name", name), slog.String("description", description), slog.String("title", title), slog.Bool("private", private))

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.CreateThread(id, userID, parentID, name, description, title, private, createdAt, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&threads.AppendCommand{
		Id:          id,
		UserId:      userID,
		ParentId:    parentID,
		Name:        name,
		Private:     private,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Title:       title,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return 
}

func (s *Controller) UpdateThread(ctx context.Context, id, userID int64, name, description, title *string) (err error) {
	slog.Debug("update thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Any("name", name), slog.Any("description", description), slog.Any("title", title))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateThread(id, userID, name, description, title, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&threads.UpdateCommand{
		Id:          id,
		UserId:      userID,
		Name:        name,
		Description: description,
		Title:       title,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.UpdateRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	slog.Debug("reorder thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.Int64("next_id", nextID), slog.Int64("prev_id", prevID))

	if parentID != -1 {
		event, err := domain.ChangeThreadParent(id, userID, parentID)
		if err != nil {
			return err
		}

		err = s.publisher.Publish(ctx, event)
		if err != nil {
			return err
		}
	}

	cmd, err := proto.Marshal(&threads.ReorderCommand{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		NextId:    nextID,
		PrevId:    prevID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.ReorderRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	slog.Debug("private thread messages", slog.Any("ids", ids), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&threads.PrivateMessagesCommand{
		Ids:    ids,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.PrivateMessagesRequest),
		Cmd:      cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	return
}

func (s *Controller) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	slog.Debug("publish thread messages", slog.Any("ids", ids), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&threads.PublishMessagesCommand{
		Ids:    ids,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.PublishMessagesRequest),
		Cmd:      cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	return
}

func (s *Controller) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("delete thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	event, err := domain.DeleteThread(id, userID)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&threads.DeleteCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.DeleteRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}
