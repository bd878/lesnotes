package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
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
	conn, err := grpc.Dial(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			cfg.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := api.NewMessagesClient(conn)

	return &Messages{cfg, client, conn}
}

func (s *Messages) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Messages) SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) (
	*model.SaveMessageResult, error,
) {
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
	}, nil
}

func (s *Messages) DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (
	*model.DeleteMessageResult, error,
) {
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

func (s *Messages) UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (
	*model.UpdateMessageResult, error,
) {
	res, err := s.client.UpdateMessage(ctx, &api.UpdateMessageRequest{
		Id: params.ID,
		UserId: params.UserID,
		FileId: params.FileID,
		Text: params.Text,
	})
	if err != nil {
		log.Errorw("client failed to save message", "error", err)
		return nil, err 
	}

	return &model.UpdateMessageResult{
		UpdateUTCNano: res.UpdateUtcNano,
	}, nil
}

func (s *Messages) ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (
	*model.ReadThreadMessagesResult, error,
) {
	log.Infow("grpc client read thread messages", "user_id", params.UserID, "thread_id", params.ThreadID)
	res, err := s.client.ReadThreadMessages(ctx, &api.ReadThreadMessagesRequest{
		UserId: params.UserID,
		ThreadId: params.ThreadID,
		Limit:  params.Limit,
		Offset: params.Offset,
		Asc:    params.Ascending,
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

func (s *Messages) ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadAllMessagesParams) (
	*model.ReadAllMessagesResult, error,
) {
	res, err := s.client.ReadAllMessages(ctx, &api.ReadAllMessagesRequest{
		UserId: int32(params.UserID),
		Limit:  params.Limit,
		Offset: params.Offset,
		Asc:    params.Ascending,
	})
	if err != nil {
		log.Errorw("client failed to read user messages", "error", err)
		return nil, err
	}

	return &model.ReadAllMessagesResult{
		Messages: model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, err
}

func (s *Messages) ReadOneMessage(ctx context.Context, log *logger.Logger, userID, messageID int32) (
	*model.Message, error,
) {
	res, err := s.client.ReadOneMessage(ctx, &api.ReadOneMessageRequest{
		UserId: userID,
		Id: messageID,
	})
	if err != nil {
		log.Errorw("client failed to read user message", "user_id", userID, "message_id", messageID, "error", err)
		return nil, err
	}

	return model.MessageFromProto(res), nil
}