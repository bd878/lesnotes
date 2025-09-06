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

type Messages struct {
	conf    Config
	client  api.MessagesClient
	conn   *grpc.ClientConn
}

func New(cfg Config) (messages *Messages) {
	messages = &Messages{conf: cfg}

	messages.setupConnection()

	return messages
}

func (s *Messages) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Messages) setupConnection() (err error) {
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

func (s *Messages) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
		return true
	}
	return false
}

func (s *Messages) SaveMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("save message", "id", id, "text", text, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private, "name", name)

	message = &model.Message{
		ID:       id,
		Text:     text,
		Name:     name,
		FileIDs:  fileIDs,
		ThreadID: threadID,
		UserID:   userID,
		Private:  private,
	}

	_, err = s.client.SaveMessage(ctx, &api.SaveMessageRequest{
		Id:       id,
		Text:     text,
		FileIds:  fileIDs,
		ThreadId: threadID,
		UserId:   userID,
		Private:  private,
		Name:     name,
	})

	return
}

func (s *Messages) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete messages", "ids", ids, "user_id", userID)

	_, err = s.client.DeleteMessages(ctx, &api.DeleteMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *Messages) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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

func (s *Messages) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
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

func (s *Messages) UpdateMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private int32) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update message", "id", id, "text", text, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private)

	_, err = s.client.UpdateMessage(ctx, &api.UpdateMessageRequest{
		Id:        id,
		UserId:    userID,
		FileIds:   fileIDs,
		Text:      text,
		Private:   private,
		ThreadId:  threadID,
	})

	return
}

func (s *Messages) ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
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
		return nil, true, err
	}

	messages = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)
	isLastPage = res.IsLastPage

	return
}

func (s *Messages) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
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

func (s *Messages) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
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
		return nil, true, err
	}

	messages = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)
	isLastPage = res.IsLastPage

	return
}

func (s *Messages) ReadMessage(ctx context.Context, id int64, userIDs []int64) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	logger.Debugw("read message", "id", id, "user_ids", userIDs)

	res, err := s.client.ReadMessage(ctx, &api.ReadMessageRequest{
		Id:      id,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	message = model.MessageFromProto(res)

	return
}

func (s *Messages) ReadPath(ctx context.Context, userID, id int64) (path []*model.Message, err error) {
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