package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/threads/pkg/loadbalance"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Config struct {
	RpcAddr  string
}

type Controller struct {
	conf         Config
	client       api.ThreadsClient
	conn         *grpc.ClientConn
}

func New(conf Config) *Controller {
	controller := &Controller{conf: conf}

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
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := api.NewThreadsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
		return true
	}
	return false
}

func (s *Controller) ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread", "id", id, "user_id", userID, "name", name)

	resp, err := s.client.Read(ctx, &api.ReadRequest{Id: id, UserId: userID, Name: name})
	if err != nil {
		return nil, err
	}

	thread = threads.ThreadFromProto(resp)

	return
}

func (s *Controller) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (list []*threads.Thread, isLastPage bool, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset, "asc", asc)

	resp, err := s.client.List(ctx, &api.ListRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
		Asc:      asc,
	})
	if err != nil {
		return nil, false, err
	}

	return threads.MapThreadsFromProto(threads.ThreadFromProto, resp.List), resp.IsLastPage, err
}

func (s *Controller) ResolveThread(ctx context.Context, id, userID int64) (path []int64, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("resolve thread", "id", id, "user_id", userID)

	resp, err := s.client.Resolve(ctx, &api.ResolveRequest{Id: id, UserId: userID})
	if err != nil {
		return nil, err
	}

	path = resp.Path

	return
}

func (s *Controller) PublishThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("publish thread", "id", id, "user_id", userID)

	_, err = s.client.Publish(ctx, &api.PublishRequest{Id: id, UserId: userID})

	return
}

func (s *Controller) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("private thread", "id", id, "user_id", userID)

	_, err = s.client.Private(ctx, &api.PrivateRequest{Id: id, UserId: userID})

	return
}


func (s *Controller) CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID,
		"next_id", nextID, "prev_id", prevID, "name", name, "description", description, "private", private)

	_, err = s.client.Create(ctx, &api.CreateRequest{
		Id:       id,
		UserId:   userID,
		ParentId: parentID,
		NextId:   nextID,
		PrevId:   prevID,
		Name:     name,
		Private:  private,
		Description: description,
	})

	return
}


func (s *Controller) UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "name", name, "description", description)

	_, err = s.client.Update(ctx, &api.UpdateRequest{
		Id:          id,
		UserId:      userID,
		Name:        name,
		Description: description,
	})

	return
}


func (s *Controller) ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("reorder thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID)

	_, err = s.client.Reorder(ctx, &api.ReorderRequest{
		Id:       id,
		UserId:   userID,
		ParentId: parentID,
		NextId:   nextID,
		PrevId:   prevID,
	})

	return
}


func (s *Controller) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete thread", "id", id, "user_id", userID)

	_, err = s.client.Delete(ctx, &api.DeleteRequest{Id: id, UserId: userID})

	return
}
