package messages

import "github.com/hashicorp/raft"

type Config struct {
  Raft raft.Config
  StreamLayer *StreamLayer
  Bootstrap bool
  DataDir string
}