package distributed

import (
  "time"
  "context"
  "bytes"

  "google.golang.org/protobuf/proto"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

func (m *DistributedMessages) apply(ctx context.Context, reqType RequestType, cmd []byte) (
  interface{}, error,
) {
  var buf bytes.Buffer
  _, err := buf.Write([]byte{byte(reqType)})
  if err != nil {
    return nil, err
  }

  _, err = buf.Write(cmd)
  if err != nil {
    return nil, err
  }

  timeout := 10*time.Second
  /* fsm.Apply() */
  future := m.raft.Apply(buf.Bytes(), timeout)
  if future.Error() != nil {
    return nil, future.Error()
  }

  res := future.Response()
  if err, ok := res.(error); ok {
    return nil, err
  }
  return res, nil
}

func (m *DistributedMessages) MakeSnapshot(ctx context.Context, log *logger.Logger) error {
  future := m.raft.Snapshot()
  if future.Error() != nil {
    return future.Error()
  }
  return nil
}

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) error {
  cmd, _ := proto.Marshal(&AppendCommand{
    Message: model.MessageToProto(params.Message),
  })

  _, err := m.apply(ctx, AppendRequest, cmd)
  if err != nil {
    log.Error("message", "raft failed to apply save message")
    return err
  }

  return nil
}

func (m *DistributedMessages) UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error {
  cmd, _ := proto.Marshal(&UpdateCommand{
    Id: params.ID,
    UserId: params.UserID,
    FileId: params.FileID,
    Text:   params.Text,
    UpdateUtcNano: params.UpdateUTCNano,
  })

  _, err := m.apply(ctx, UpdateRequest, cmd)
  if err != nil {
    log.Error("message", "raft failed to apply save message")
    return err
  }

  return nil
}

func (m *DistributedMessages) DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
  cmd, _ := proto.Marshal(&DeleteCommand{
    Id: params.ID,
    UserId: params.UserID,
    FileId: params.FileID,
  })

  _, err := m.apply(ctx, DeleteRequest, cmd)
  if err != nil {
    log.Error("message", "raft failed to apply delete message")
    return err
  }

  return nil
}

func (m *DistributedMessages) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.ReadUserMessagesResult, error,
) {
  return m.repo.ReadUserMessages(
    ctx,
    log,
    &model.ReadUserMessagesParams{
      UserID:    params.UserID,
      Limit:     params.Limit,
      Offset:    params.Offset,
      Ascending: params.Ascending,
    },
  )
}
