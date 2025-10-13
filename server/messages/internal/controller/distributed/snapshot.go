package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"github.com/bd878/gallery/server/logger"
)

type snapshot struct {
	reader io.ReadCloser
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	reader, err := f.repo.Dump(context.Background())
	if err != nil {
		if err2 := reader.Close(); err2 != nil {
			logger.Errorw("failed to close snapshot reader", "error", err2)
		}

		return nil, err
	}

	return &snapshot{reader}, nil
}

func (f *fsm) Restore(reader io.ReadCloser) (err error) {
	err = f.repo.Truncate(context.Background())
	if err != nil {
		logger.Errorw("truncate returned error", "error", err)
	}
	defer reader.Close()
	return f.repo.Restore(context.Background(), reader)
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	_, err = io.Copy(sink, s.reader)
	defer sink.Cancel()
	if err != nil {
		return
	}

	return sink.Close()
}

func (s *snapshot) Release() {
	if err := s.reader.Close(); err != nil {
		logger.Errorw("failed to close snapshot reader", "error", err)
	}
}
