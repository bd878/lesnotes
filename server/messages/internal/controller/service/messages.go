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
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Config struct {
	RpcAddr string
}

type ThreadsGateway interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32) (list []*threads.Thread, isLastPage bool, err error)
	// TODO: fix params order
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	CreateThread(ctx context.Context, id, userID, parentID int64, name string, private bool) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	UpdateThread(ctx context.Context, id, userID int64) (err error)
	ResolvePath(ctx context.Context, userID, id int64) (path []int64, err error)
	ReadThread(ctx context.Context, userID, id int64) (thread *threads.Thread, err error)
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

func (s *Controller) UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "user_id", userID)

	err = s.threads.UpdateThread(ctx, id, userID)
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
	})

	return
}

// Get messages in order
func (s *Controller) ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (list *model.List, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread messages", "user_id", userID, "thread_id", threadID, "limit", limit, "offset", offset, "ascending", ascending)

	total, err := s.threads.CountThreads(ctx, threadID, userID)
	if err != nil {
		return nil, err
	}

	threadsList, isLastPage, err := s.threads.ListThreads(ctx, userID, threadID, limit, offset)
	if err != nil {
		logger.Debugw("failed to list threads", "error", err)
		return nil, err
	}

	ids := make([]int64, 0)
	for _, thread := range threadsList {
		ids = append(ids, thread.ID)
	}

	res, err := s.client.ReadBatchMessages(ctx, &api.ReadBatchMessagesRequest{
		UserId:   userID,
		Ids:      ids,
	})
	if err != nil {
		logger.Debugw("failed to read batch messages", "error", err)
		return nil, err
	}

	messages := model.MapMessagesFromProto(model.MessageFromProto, res.Messages)
	for i, message := range messages {
		message.Count = threadsList[i].Count
	}

	list = &model.List{
		Messages:      messages,
		IsLastPage:    isLastPage,
		IsFirstPage:   offset == 0,
		Total:         total,
		Count:         int32(len(ids)),
		Offset:        offset,
	}

	return
}

// Read messages by given ids
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

// read all messages not concerning thread
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

	list = &model.List{
		Messages:     model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage:   res.IsLastPage,
		IsFirstPage:  offset == 0,
		Offset:       offset,
		// TODO: Total:        total.Count threads.CountThreads,
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

func (s *Controller) ReadPath(ctx context.Context, userID, id int64) (path []*model.Message, parentID int64, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, 0, err
		}
	}

	logger.Debugw("read path", "user_id", userID, "id", id)

	ids, err := s.threads.ResolvePath(ctx, userID, id)
	if err != nil {
		return nil, 0, err
	}

	// TODO: log error, falls if error 0
	thread, _ := s.threads.ReadThread(ctx, userID, id)
	if thread != nil {
		parentID = thread.ParentID
	}

	res, err := s.client.ReadBatchMessages(ctx, &api.ReadBatchMessagesRequest{
		UserId: userID,
		Ids:    ids,
	})
	if err != nil {
		return nil, 0, err
	}

	path = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)

	return
}
