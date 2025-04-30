package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/internal/loadbalance"
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

func New(cfg Config) *Messages {
	msg := &Messages{conf: cfg}
	msg.setupConnection()
	return msg
}

func (s *Messages) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Messages) setupConnection() error {
	conn, err := grpc.Dial(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("messages cannot setup grpc client connection, exit")
		return err
	}

	client := api.NewMessagesClient(conn)

	s.conn = conn
	s.client = client
	return nil
}

func (s *Messages) isConnFailed() bool {
	state := s.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (s *Messages) SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) (
	*model.SaveMessageResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	res, err := s.client.SaveMessage(ctx, &api.SaveMessageRequest{
		Message: model.MessageToProto(message),
	})
	if err != nil {
		log.Errorw("client failed to save message", "error", err)
		return nil, err 
	}

	return &model.SaveMessageResult{
		ID: res.Id,
		UpdateUTCNano: res.UpdateUtcNano,
		CreateUTCNano: res.CreateUtcNano,
		Private: res.Private,
	}, nil
}

func (s *Messages) DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (
	*model.DeleteMessageResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	_, err := s.client.DeleteMessage(ctx, &api.DeleteMessageRequest{
		Id: params.ID,
		UserId: params.UserID,
	})
	if err != nil {
		log.Errorw("client failed to delete message", "error", err)
		return nil, err
	}

	return &model.DeleteMessageResult{}, nil
}

func (s *Messages) DeleteMessages(ctx context.Context, log *logger.Logger, params *model.DeleteMessagesParams) (
	*model.DeleteMessagesResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	res, err := s.client.DeleteMessages(ctx, &api.DeleteMessagesRequest{
		Ids: params.IDs,
		UserId: params.UserID,
	})
	if err != nil {
		log.Errorw("client failed to delete batch messages", "error", err)
		return nil, err
	}

	var ids []*model.DeleteMessageStatus
	for _, status := range res.Ids {
		ids = append(ids, &model.DeleteMessageStatus{
			OK: status.Ok,
			Explain: status.Explain,
		})
	}

	return &model.DeleteMessagesResult{
		IDs: ids,
	}, nil
}

func (s *Messages) UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (
	*model.UpdateMessageResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	resp, err := s.client.UpdateMessage(ctx, &api.UpdateMessageRequest{
		Id: params.ID,
		UserId: params.UserID,
		FileId: params.FileID,
		Text: params.Text,
		Private: params.Private,
	})
	if err != nil {
		log.Errorw("client failed to save message", "error", err)
		return nil, err 
	}

	return &model.UpdateMessageResult{
		ID: params.ID,
		UpdateUTCNano: resp.UpdateUtcNano,
		Private: resp.Private,
	}, nil
}

func (s *Messages) ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (
	*model.ReadThreadMessagesResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	log.Infow("grpc client read thread messages", "user_id", params.UserID, "thread_id", params.ThreadID, "private", params.Private)
	res, err := s.client.ReadThreadMessages(ctx, &api.ReadThreadMessagesRequest{
		UserId: params.UserID,
		ThreadId: params.ThreadID,
		Limit:  params.Limit,
		Offset: params.Offset,
		Asc:    params.Ascending,
		Private: params.Private,
	})
	if err != nil {
		log.Errorw("client failed to read thread messages", "thread_id", params.ThreadID)
		return nil, err
	}

	return &model.ReadThreadMessagesResult{
		Messages: model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, err
}

func (s *Messages) ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (
	*model.ReadMessagesResult, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	res, err := s.client.ReadAllMessages(ctx, &api.ReadMessagesRequest{
		UserId: int32(params.UserID),
		Limit:  params.Limit,
		Offset: params.Offset,
		Asc:    params.Ascending,
		Private: params.Private,
	})
	if err != nil {
		log.Errorw("client failed to read user messages", "error", err)
		return nil, err
	}

	return &model.ReadMessagesResult{
		Messages: model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, err
}

func (s *Messages) ReadOneMessage(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (
	*model.Message, error,
) {
	if s.isConnFailed() {
		log.Info("conn failed, setup new connection")
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	res, err := s.client.ReadOneMessage(ctx, &api.ReadOneMessageRequest{
		Id: params.ID,
		UserIds: params.UserIDs,
	})
	if err != nil {
		log.Errorw("client failed to read user message", "user_ids", params.UserIDs, "message_id", params.ID, "error", err)
		return nil, err
	}

	return model.MessageFromProto(res), nil
}