package distributed

import (
  "io"
  "context"
  "errors"
  "encoding/json"
  "bytes"

  "github.com/hashicorp/raft"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/internal/repository"
)

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
  var msg *model.Message
  var err error

  msg, err = f.repo.FindByIndexTerm(context.Background(), logger.Default(), &model.FindByIndexParams{
    LogIndex: record.Index,
    LogTerm:  record.Term,
  })
  if err != nil {
    /* not found is expected behaviour */
    if !errors.Is(err, repository.ErrNotFound) {
      return err
    }
  }
  if msg != nil {
    return ErrMsgExist
  }

  buf := record.Data
  err = json.Unmarshal(buf, &msg)
  if err != nil {
    return err
  }
  msg.LogIndex = record.Index
  msg.LogTerm = record.Term

  msg.Id, err = f.repo.Put(context.Background(), logger.Default(), &model.PutParams{Message: msg})
  if err != nil {
    return err
  }

  return *msg
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
  return &snapshot{repo: f.repo}, nil
}

// TODO: restore will reapply same messages with ids,
// check whether msg with id exists
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