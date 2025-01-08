package distributed

import (
  "time"
  "context"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (
  resMsg *model.Message, err error,
) {
  params.Message.CreateTime = time.Now().String()
  if resMsg, err = m.apply(ctx, params.Message); err != nil {
    log.Error("message", "raft failed to apply save message")
    return nil, err
  }
  return resMsg, nil
}

func (m *DistributedMessages) UpdateMessage(_ context.Context, _ *logger.Logger, _ *model.UpdateMessageParams) (
  *model.Message, error,
) {
  /* not implemented */
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
