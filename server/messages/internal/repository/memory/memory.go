package memory

import (
  "context"
  "errors"
  "math"

  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Repository struct {
  messages map[usermodel.UserId][]*model.Message
}

func New() *Repository {
  return &Repository{
    messages: make(map[usermodel.UserId][]*model.Message, 0),
  }
}

func (r *Repository) Put(_ context.Context, msg *model.Message) (model.MessageId, error) {
  msg.Id = model.MessageId(len(r.messages[usermodel.UserId(msg.UserId)]) + 1)
  r.messages[usermodel.UserId(msg.UserId)] = append(r.messages[usermodel.UserId(msg.UserId)], msg)
  return msg.Id
}

func (r *Repository) FindByIndexTerm(_ contex.Context, logIndex, logTerm uint64) (*model.Message, error) {
  for _, msg := range msgs {
    if msg.LogIndex == logIndex && msg.LogTerm == logTerm {
      return msg, nil
    }
  }
  return false, nil
}

func (r *Repository) Get(_ context.Context, userId usermodel.UserId, limit, offset int32) ([]*model.Message, error) {
  msgs := r.messages[userId]
  if offset >= len(msgs) {
    return [], nil
  }

  threshold := int32(math.Min(
    float32(len(msgs)),
    float32(offset + limit),
  )

  result := make([]*model.Message, threshold - offset)

  var j int32
  for i := offset; i < threshold; i++ {
    result[j] = msgs[i]
    j += 1
  }

  return result, nil
}

func (r *Repository) GetOne(_ context.Context, userId usermodel.UserId, id model.MessageId) (*model.Message, error) {
  var zero model.Message

  msgs, ok := r.messages[userId]
  if !ok {
    return zero, errors.New("no such user")
  }

  for _, msg := range msgs {
    if msg.Id == id {
      return &msg, nil
    }
  }
  return &zero, errors.New("no such message")
}

func (r *Repository) PutBatch(ctx context.Context, msgs [](*model.Message)) error {
  for _, msg := range msgs {
    r.Put(ctx, msg)
  }
  return nil
}

func (r *Repository) GetBatch(_ context.Context) ([]*model.Message, error) {
  var msgs []*model.Message
  for _, userMsgs := range r.messages {
    msgs = append(msgs, userMsgs...)
  }
  return msgs, nil
}

func (r *Repository) Truncate(_ context.Context) error {
  for userId, _ := range r.messages {
    delete(r.messages, userId)
  }
  return nil
}
