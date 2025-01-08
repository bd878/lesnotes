package distributed

import (
  "io"
  "context"
  "encoding/json"
  "bytes"

  "github.com/hashicorp/raft"
  "github.com/bd878/gallery/server/logger"
)

type snapshot struct {
  repo Repository
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
  msgs, err := s.repo.GetBatch(context.Background(), logger.Default())
  if err != nil {
    return err
  }

  b, err := json.Marshal(msgs)
  if err != nil {
    return err
  }
  if _, err := io.Copy(sink, bytes.NewReader(b)); err != nil {
    _ = sink.Cancel()
    return err
  }
  return sink.Close()
}

func (s *snapshot) Release() {}
