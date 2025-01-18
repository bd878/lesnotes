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

func (s *Messages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (
  *model.SaveMessageResult, error,
) {
  res, err := s.client.SaveMessage(ctx, &api.SaveMessageRequest{
    Message: model.MessageToProto(params.Message),
  })
  if err != nil {
    log.Error("message", "client failed to save message")
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
    log.Error("message", "client failed to delete message")
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
    log.Error("message", "client failed to save message")
    return nil, err 
  }

  return &model.UpdateMessageResult{
    UpdateUTCNano: res.UpdateUtcNano,
  }, nil
}

func (s *Messages) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.ReadUserMessagesResult, error,
) {
  res, err := s.client.ReadUserMessages(ctx, &api.ReadUserMessagesRequest{
    UserId: int32(params.UserID),
    Limit:  params.Limit,
    Offset: params.Offset,
    Asc:    params.Ascending,
  })
  if err != nil {
    log.Error("message", "client failed to read user messages")
    return nil, err
  }

  return &model.ReadUserMessagesResult{
    Messages: model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
    IsLastPage: res.IsLastPage,
  }, err
}
