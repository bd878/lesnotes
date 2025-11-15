package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type Config struct {
	RpcAddr string
}

type ThreadsGateway interface {
	CreateThread(ctx context.Context, id, userID, parentID int64, name string, private bool) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	UpdateThread(ctx context.Context, id, userID, parentID int64) (err error)
}

type Controller struct {
	conf       Config
	client     api.MessagesClient
	conn       *grpc.ClientConn
	threads    ThreadsGateway
}

func New(conf Config, threads ThreadsGateway) *Controller {
	controller := &Controller{conf: conf, threads: threads}

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

	client := api.NewMessagesClient(conn)

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

func (s *Controller) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("save message", "id", id, "text", text, "title", title, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private, "name", name)

	err = s.threads.CreateThread(ctx, id, userID, threadID, name, private)
	if err != nil {
		return
	}

	_, err = s.client.SaveMessage(ctx, &api.SaveMessageRequest{
		Id:       id,
		Text:     text,
		Title:    title,
		FileIds:  fileIDs,
		ThreadId: threadID,
		UserId:   userID,
		Private:  private,
		Name:     name,
	})
	if err != nil {
		return
	}

	message = &model.Message{
		ID:       id,
		Text:     text,
		Title:    title,
		Name:     name,
		FileIDs:  fileIDs,
		ThreadID: threadID,
		UserID:   userID,
		Private:  private,
	}

	return
}

func (s *Controller) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete messages", "ids", ids, "user_id", userID)

	// TODO: DeleteThreads
	for _, id := range ids {
		err = s.threads.DeleteThread(ctx, id, userID)
		if err != nil {
			return
		}
	}

	_, err = s.client.DeleteMessages(ctx, &api.DeleteMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *Controller) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("publish messages", "ids", ids, "user_id", userID)

	_, err = s.client.PublishMessages(ctx, &api.PublishMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *Controller) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("private messages", "ids", ids, "user_id", userID)

	_, err = s.client.PrivateMessages(ctx, &api.PrivateMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *Controller) UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, threadID int64, userID int64, private int32) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private)

	err = s.threads.UpdateThread(ctx, id, userID, threadID)
	if err != nil {
		return
	}

	_, err = s.client.UpdateMessage(ctx, &api.UpdateMessageRequest{
		Id:        id,
		UserId:    userID,
		FileIds:   fileIDs,
		Text:      text,
		Title:     title,
		Name:      name,
		Private:   private,
		ThreadId:  threadID,
	})

	return
}

func (s *Controller) ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (list *model.List, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread messages", "user_id", userID, "thread_id", threadID, "limit", limit, "offset", offset, "ascending", ascending)

	res, err := s.client.ReadThreadMessages(ctx, &api.ReadThreadMessagesRequest{
		UserId:   userID,
		ThreadId: threadID,
		Limit:    limit,
		Offset:   offset,
		Asc:      ascending,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.client.CountMessages(ctx, &api.CountMessagesRequest{
		UserId:   userID,
		ThreadId: threadID,
	})
	if err != nil {
		return nil, err
	}

	list = &model.List{
		Messages:      model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage:    res.IsLastPage,
		IsFirstPage:   offset == 0,
		Total:         total.Count,
		Count:         int32(len(res.Messages)),
		Offset:        offset,
	}

	return
}

func (s *Controller) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read batch messages", "user_id", userID, "ids", ids)

	res, err := s.client.ReadBatchMessages(ctx, &api.ReadBatchMessagesRequest{
		UserId:   userID,
		Ids:      ids,
	})
	if err != nil {
		return nil, err
	}

	messages = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)

	return
}

func (s *Controller) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (list *model.List, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)

	res, err := s.client.ReadMessages(ctx, &api.ReadMessagesRequest{
		UserId:   userID,
		Limit:    limit,
		Offset:   offset,
		Asc:      ascending,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.client.CountMessages(ctx, &api.CountMessagesRequest{
		UserId:   userID,
		ThreadId: -1,
	})
	if err != nil {
		return nil, err
	}

	list = &model.List{
		Messages:     model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage:   res.IsLastPage,
		IsFirstPage:  offset == 0,
		Offset:       offset,
		Total:        total.Count,
		Count:        int32(len(res.Messages)),
	}

	return
}

func (s *Controller) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	logger.Debugw("read message", "id", id, "name", name, "user_ids", userIDs)

	res, err := s.client.ReadMessage(ctx, &api.ReadMessageRequest{
		Id:      id,
		UserIds: userIDs,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}

	message = model.MessageFromProto(res)

	return
}

func (s *Controller) ReadPath(ctx context.Context, userID, id int64) (path []*model.Message, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	logger.Debugw("read path", "user_id", userID, "id", id)

	res, err := s.client.ReadPath(ctx, &api.ReadPathRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	path = model.MapMessagesFromProto(model.MessageFromProto, res.Path)

	return
}

func (s *Controller) ReadMessagesAround(ctx context.Context, userID, threadID, id int64, limit int32) (list *model.List, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read messages around", "user_id", userID, "thread_id", threadID, "id", id, "limit", limit)

	res, err := s.client.ReadMessagesAround(ctx, &api.ReadMessagesAroundRequest{
		UserId:   userID,
		ThreadId: threadID,
		Limit:    limit,
		Id:       id,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.client.CountMessages(ctx, &api.CountMessagesRequest{
		UserId:   userID,
		ThreadId: threadID,
	})
	if err != nil {
		return nil, err
	}

	list = &model.List{
		Messages:    model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage:  res.IsLastPage,
		IsFirstPage: res.Offset == 0,
		Offset:      res.Offset,
		Count:       int32(len(res.Messages)),
		Total:       total.Count,
	}

	return
}
