package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/threads/internal/domain"
	"github.com/bd878/gallery/server/db/threads/pkg/machine"
	"github.com/bd878/gallery/server/db/threads/pkg/loadbalance"
	"github.com/bd878/gallery/server/threads/pkg/model"
	"github.com/bd878/gallery/server/threads/config"
)

type Controller struct {
	conf         config.Config
	client       threads.ThreadsClient
	conn         *grpc.ClientConn
	publisher    ddd.EventPublisher[ddd.Event]
}

func New(conf config.Config, publisher ddd.EventPublisher[ddd.Event]) *Controller {
	controller := &Controller{
		conf: conf,
		publisher: publisher,
	}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Controller) setupConnection() (err error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.ThreadsServiceAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := threads.NewThreadsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		logger.Debugw("threads conn failed", "state", state.String())
		return true
	}
	return false
}

func (s *Controller) ReadThread(ctx context.Context, id, userID int64, name string) (thread *model.Thread, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread", "id", id, "user_id", userID, "name", name)

	resp, err := s.client.Read(ctx, &threads.ReadRequest{Id: id, UserId: userID, Name: name})
	if err != nil {
		return nil, err
	}

	thread = model.ThreadFromProto(resp)

	return
}

func (s *Controller) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*model.Thread, isLastPage bool, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc)

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
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list messages", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc, "private_message", privateMessage)

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
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("resolve thread", "id", id, "user_id", userID)

	resp, err := s.client.Resolve(ctx, &threads.ResolveRequest{Id: id, UserId: userID})
	if err != nil {
		return nil, err
	}

	path = model.PathFromProto(resp.Path)

	return
}

func (s *Controller) PublishThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("publish thread", "id", id, "user_id", userID)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&threads.PublishCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	event, err := domain.PublishThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *Controller) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("private thread", "id", id, "user_id", userID)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&threads.PrivateCommand{
		Id:            id,
		UserId:        userID,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	event, err := domain.PrivateThread(id, userID, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}


func (s *Controller) CreateThread(ctx context.Context, id, userID, parentID int64, name, description, title string, private bool) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID,
		"name", name, "description", description, "title", title, "private", private)


	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.CreateThread(id, userID, parentID, name, description, title, private, createdAt, updatedAt)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&threads.AppendCommand{
		Id:            id,
		UserId:        userID,
		ParentId:      parentID,
		Name:          name,
		Private:       private,
		Description:   description,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Title:         title,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	return s.publisher.Publish(ctx, event)
}

func (s *Controller) UpdateThread(ctx context.Context, id, userID int64, name, description, title *string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "name", name, "description", description, "title", title)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateThread(id, userID, name, description, title, updatedAt)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&threads.UpdateCommand{
		Id:            id,
		UserId:        userID,
		Name:          name,
		Description:   description,
		Title:         title,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	return s.publisher.Publish(ctx, event)
}


func (s *Controller) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("reorder thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID)

	cmd, err := proto.Marshal(&threads.ReorderCommand{
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

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.ReorderRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	if parentID != -1 {
		event, err := domain.ChangeThreadParent(id, userID, parentID)
		if err != nil {
			return err
		}

		err = s.publisher.Publish(ctx, event)
	}

	return
}


func (s *Controller) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("private thread messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&threads.PrivateMessagesCommand{
		Ids:         ids,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateMessagesRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	return
}

func (s *Controller) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("publish thread messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&threads.PublishMessagesCommand{
		Ids:       ids,
		UserId:    userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishMessagesRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	return
}

func (s *Controller) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&threads.DeleteCommand{
		Id:       id,
		UserId:   userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	event, err := domain.DeleteThread(id, userID)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}
