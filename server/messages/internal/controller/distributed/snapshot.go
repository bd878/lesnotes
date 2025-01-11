package distributed

import (
  "github.com/hashicorp/raft"
)

type snapshot struct {
  repo Repository
}

func (s *snapshot) Persist(_ raft.SnapshotSink) error {
  // TODO: .dump entier database from fsm.Restore()
  return nil
}

func (s *snapshot) Release() {}
