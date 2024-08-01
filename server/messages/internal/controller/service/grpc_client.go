package service

import (
  "context"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Config struct {
  ClusterNodeAddr string
}

type Messages struct {
  cfg    Config
  client api.MessagesClient
  conn   *grpc.ClientConn
}

func New(cfg Config) *Messages {
  conn, err := grpc.Dial(cfg.ClusterNodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (s *Messages) SaveMessage(ctx context.Context, msg *api.Message) (
  model.MessageId,
  error,
) {
  res, err := s.client.SaveMessage(ctx, &api.SaveMessageRequest{Message: msg})
  if err != nil {
    return model.MessageId(-1), err 
  }
  return model.MessageId(res.Id), nil
}

func (s *Messages) ReadUserMessages(ctx context.Context, userId usermodel.UserId) (
  []*api.Message,
  error,
) {
  res, err := s.client.ReadUserMessages(ctx, &api.ReadUserMessagesRequest{UserId: uint32(userId)})
  if err != nil {
    return nil, err
  }
  return res.Messages, err
}