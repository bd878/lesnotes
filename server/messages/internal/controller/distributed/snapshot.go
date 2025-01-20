package distributed

import (
  "github.com/hashicorp/raft"
)

type snapshot struct {}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
  if _, err := sink.Write([]byte{}); err != nil {
    _ = sink.Cancel()
    return err
  }
  return sink.Close()
}

func (s *snapshot) Release() {}
