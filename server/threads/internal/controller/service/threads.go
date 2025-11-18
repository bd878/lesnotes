package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
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

func (s *Controller) ReadThread(ctx context.Context, id, userID int64) (thread *threads.Thread, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread", "id", id, "user_id", userID)

	resp, err := s.client.Read(ctx, &api.ReadRequest{Id: id, UserId: userID})
	if err != nil {
		return nil, err
	}

	thread = threads.ThreadFromProto(resp)

	return
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


func (s *Controller) CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name string, private bool) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID, "next_id", nextID, "prev_id", prevID, "name", name, "private", private)

	_, err = s.client.Create(ctx, &api.CreateRequest{
		Id:       id,
		UserId:   userID,
		ParentId: parentID,
		NextId:   nextID,
		PrevId:   prevID,
		Name:     name,
		Private:  private,
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
