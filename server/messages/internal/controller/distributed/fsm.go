package distributed

import (
  "io"
  "context"

  "github.com/hashicorp/raft"
  "google.golang.org/protobuf/proto"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/logger"
)

type Repository interface {
  Put(ctx context.Context, log *logger.Logger, params *model.PutParams) (int32, error)
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
  // Put does not put message with same id twice
  res, err := f.repo.Put(context.Background(), logger.Default(), &model.PutParams{
    Message: model.MessageFromProto(cmd.Message),
  })
  if err != nil {
    return nil
  }
  return &AppendCommandResult{
    Id: res,
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

// restore from snapshot (db .dump)
func (f *fsm) Restore(_ io.ReadCloser) error {
  /* TODO: make repository dump */
  return nil
}