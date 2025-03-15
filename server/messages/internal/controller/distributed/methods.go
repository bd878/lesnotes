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

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) error {
  cmd, _ := proto.Marshal(&AppendCommand{
    Message: model.MessageToProto(message),
  })

  _, err := m.apply(ctx, AppendRequest, cmd)
  if err != nil {
    log.Errorw("raft failed to apply save message", "error", err)
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
    log.Errorw("raft failed to apply save message", "error", err)
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
    log.Errorw("raft failed to apply delete message", "error", err)
    return err
  }

  return nil
}

func (m *DistributedMessages) ReadMessage(ctx context.Context, log *logger.Logger, userID, messageID int32) (
  *model.Message, error,
) {
  return m.repo.Read(ctx, log, messageID)
}

func (m *DistributedMessages) ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (
  *model.ReadThreadMessagesResult, error,
) {
  log.Infow("distributed methods. read thread messages", "user_id", params.UserID, "thread_id", params.ThreadID)
  return m.repo.ReadThreadMessages(
    ctx,
    log,
    &model.ReadThreadMessagesParams{
      UserID:    params.UserID,
      ThreadID:  params.ThreadID,
      Limit:     params.Limit,
      Offset:    params.Offset,
      Ascending: params.Ascending,
    },
  )
}

func (m *DistributedMessages) ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadAllMessagesParams) (
  *model.ReadAllMessagesResult, error,
) {
  return m.repo.ReadAllMessages(
    ctx,
    log,
    &model.ReadAllMessagesParams{
      UserID:    params.UserID,
      Limit:     params.Limit,
      Offset:    params.Offset,
      Ascending: params.Ascending,
    },
  )
}
