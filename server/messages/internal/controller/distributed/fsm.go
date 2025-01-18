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
  Create(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) error
  Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error
  Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error
  ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (*model.ReadUserMessagesResult, error)
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
  proto.Unmarshal(raw, &cmd)

  // Put does not put message with same id twice
  err := f.repo.Create(context.Background(), logger.Default(), &model.SaveMessageParams{
    Message: model.MessageFromProto(cmd.Message),
  })
  if err != nil {
    return err
  }
  return &AppendCommandResult{}
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
  var cmd UpdateCommand
  proto.Unmarshal(raw, &cmd)

  err := f.repo.Update(context.Background(), logger.Default(), &model.UpdateMessageParams{
    ID: cmd.Id,
    UserID: cmd.UserId,
    FileID: cmd.FileId,
    Text: cmd.Text,
    UpdateUTCNano: cmd.UpdateUtcNano,
  })
  if err != nil {
    return err
  }

  return &UpdateCommandResult{
  }
}

func (f *fsm) applyDelete(raw []byte) interface{} {
  var cmd DeleteCommand
  proto.Unmarshal(raw, &cmd)

  err := f.repo.Delete(context.Background(), logger.Default(), &model.DeleteMessageParams{
    ID: cmd.Id,
    UserID: cmd.UserId,
    FileID: cmd.FileId,
  })
  if err != nil {
    return err
  }

  return &DeleteCommandResult{}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
  return &snapshot{repo: f.repo}, nil
}

// restore from snapshot (db .dump)
func (f *fsm) Restore(_ io.ReadCloser) error {
  /* TODO: make repository dump */
  return nil
}