package distributed

import (
  "io"
  "context"
  "encoding/json"
  "bytes"

  "github.com/hashicorp/raft"
  "google.golang.org/protobuf/proto"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/logger"
)

type Repository interface {
  Put(ctx context.Context, log *logger.Logger, params *model.PutParams) (model.MessageId, error)
  Get(ctx context.Context, log *logger.Logger, params *model.GetParams) (*model.MessagesList, error)
  GetBatch(ctx context.Context, log *logger.Logger) ([]*model.Message, error)
  Truncate(ctx context.Context, log *logger.Logger) error
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
  repo Repository
}

/**
 * Returns empty interface. It is either an error,
 * or new msg with unique id, saved in repo.
 * 
 * Apply replicates log state from the bottom up.
 * Leader makes Apply on start.
 */
func (f *fsm) Apply(record *raft.Log) interface{} {
  buf := record.Data
  reqType := RequestType(buf[0])
  switch reqType {
  case AppendRequest:
    return f.applyAppend(buf[1:])
  case UpdateRequest:
    return f.applyUpdate(buf[1:])
  case DeleteRequest:
    return f.applyDelete(buf[1:])
  default:
    logger.Error("unknown request type: ", reqType)
  }
  return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
  var cmd AppendCommand
  err := proto.Unmarshal(raw, &cmd)
  if err != nil {
    return err
  }
  res, err := f.repo.Put(context.Background(), logger.Default(), &model.PutParams{
    Message: model.MessageFromProto(cmd.Message),
  })
  if err != nil {
    return nil
  }
  return &AppendCommandResult{
    Id: uint32(res),
  }
}

func (f *fsm) applyUpdate(_ []byte) interface{} {
  /* not implemented */
  return nil
}

func (f *fsm) applyDelete(_ []byte) interface{} {
  /* not implemented */
  return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
  return &snapshot{repo: f.repo}, nil
}

// restore will reapply same messages with ids
func (f *fsm) Restore(r io.ReadCloser) error {
  var buf *bytes.Buffer
  var msgs []model.Message

  _, err := io.Copy(buf, r)
  if err == io.EOF {
    return err
  } else if err != nil {
    return err
  }
  err = json.Unmarshal(buf.Bytes(), &msgs)
  if err != nil {
    return err
  }

  ctx := context.Background()
  err = f.repo.Truncate(ctx, logger.Default())
  if err != nil {
    return err
  }
  for _, msg := range msgs {
    _, err := f.repo.Put(ctx, logger.Default(), &model.PutParams{Message: &msg})
    if err != nil {
      return err
    }
  }
  return nil
}