package memory

import (
  "context"
  "errors"

  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/user/pkg/model"
)

type Repository struct {
  messages map[usermodel.UserId][]model.Message
}

func New() *Repository {
  return &Repository{
    messages: make(map[usermodel.UserId][]model.Message, 0),
  }
}

func (r *Repository) Put(_ context.Context, msg *model.Message) error {
  r.messages[usermodel.UserId(msg.UserId)] = append(r.messages[usermodel.UserId(msg.UserId)], *msg)
  return nil
}

func (r *Repository) Get(_ context.Context, userId usermodel.UserId) ([]model.Message, error) {
  return r.messages[userId], nil
}

func (r *Repository) GetOne(_ context.Context, userId usermodel.UserId, id int) (model.Message, error) {
  var zero model.Message

  msgs, ok := r.messages[userId]
  if !ok {
    return zero, errors.New("no such user")
  }

  for _, msg := range msgs {
    if msg.Id == id {
      return msg, nil
    }
  }
  return zero, errors.New("no such message")
}

func (r *Repository) PutBatch(ctx context.Context, msgs [](*model.Message)) error {
  for _, msg := range msgs {
    r.Put(ctx, msg)
  }
  return nil
}

func (r *Repository) GetBatch(_ context.Context) ([]model.Message, error) {
  var msgs []model.Message
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
