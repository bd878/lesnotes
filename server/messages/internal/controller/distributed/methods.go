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

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (
  *model.Message, error,
) {
  cmd, err := proto.Marshal(&AppendCommand{
    Message: model.MessageToProto(params.Message),
  })
  if err != nil {
    return nil, err
  }

  res, err := m.apply(ctx, AppendRequest, cmd)
  if err != nil {
    log.Error("message", "raft failed to apply save message")
    return nil, err
  }

  params.Message.Id = model.MessageId(res.(*AppendCommandResult).Id)

  return params.Message, nil
}

func (m *DistributedMessages) UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (
  *model.Message, error,
) {
  cmd, err := proto.Marshal(&UpdateCommand{
    Message: model.MessageToProto(params.Message),
  })
  if err != nil {
    return nil, err
  }

  _, err = m.apply(ctx, UpdateRequest, cmd)
  if err != nil {
    log.Error("message", "raft failed to apply save message")
    return nil, err
  }

  /* TODO: map res interface{} to model.Message */

  return nil, nil
}

func (m *DistributedMessages) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.MessagesList, error,
) {
  return m.repo.Get(
    ctx,
    log,
    &model.GetParams{
      UserId:    params.UserId,
      Limit:     params.Limit,
      Offset:    params.Offset,
      Ascending: params.Ascending,
    },
  )
}
