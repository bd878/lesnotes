package machine

import (
	"io"
	"os"
	"context"

	"go.uber.org/zap"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/store"
)

type snapshot struct {
	store         *store.Store
	dumper   Dumper
	ctx           context.Context
	ch            <-chan *threads.ThreadsSnapshot
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	logger.Debugln("snapshotting threads")

	s := &snapshot{}

	file, err := os.CreateTemp("", "threads_*.bin")
	if err != nil {
		return nil, err
	}

	s.store, err = store.NewStore(file)
	if err != nil {
		return nil, err
	}

	s.dumper = f.dumper
	s.ctx = context.TODO()
	s.ch, err = f.dumper.Open(s.ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// raft NewFileSnapshotStore has a bug : bufio.Reader has 4096 default size.
// When reading per bytes, it may occasionally return 0 bytes off limit.
// With bufio, you must explicitely set larger buffer amount (file size)
func (f *Machine) Restore(reader io.ReadCloser) (err error) {
	logger.Debugln("restoring fsm from snapshot")

	s := store.NewReader(reader)
	defer s.Close()

	for {
		size, err := s.ReadSize()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		logger.Debugw("restore", "size", size)

		data := make([]byte, size)
		n, err := s.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		logger.Debugw("restore", "n", n)

		if uint64(n) < size {
			n2, err := s.Read(data[n:])
			if err != nil {
				return err
			}

			logger.Debugw("restore", "n2", n2)
		}

		var snapshot threads.ThreadsSnapshot
		if err = proto.Unmarshal(data, &snapshot); err != nil {
			logger.Debugln(zap.Error(err))
			return err
		}

		bytes, err := prototext.Marshal(&snapshot)
		if err != nil {
			logger.Debugln(zap.Error(err))
			continue
		}

		logger.Debugln(zap.String("snapshot", string(bytes)))

		err = f.dumper.Restore(context.TODO(), &snapshot)
		if err != nil {
			return err
		}
	}

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	logger.Debugln("persisting snapshot")

	for snapshot := range s.ch {
		switch v := snapshot.Item.(type) {
		case *threads.ThreadsSnapshot_Thread:
			logger.Debugw("thread snapshot", "id", v.Thread.Id)
		default:
			logger.Debugln("unknown snapshot")
			continue
		}

		data, err := proto.Marshal(snapshot)
		if err != nil {
			return err
		}

		n, err := s.store.Append(data)
		if err != nil {
			return err
		}

		logger.Debugw("persist", "n", n)

		select {
		case <-s.ctx.Done():
			return context.Cause(s.ctx)
		default:
		}
	}

	logger.Debugln("seek store")

	err = s.store.Seek()
	if err != nil {
		return
	}

	n, err := io.Copy(sink, s.store.File)
	if err != nil {
		return err
	}

	logger.Debugw("store persisted", "n", n)

	return
}

func (s *snapshot) Release() {
	logger.Debugln("release snapshot")
	if err := s.store.Close(); err != nil {
		logger.Errorw("cannot close store file", "error", err)
	}

	if err := s.dumper.Close(); err != nil {
		logger.Errorw("cannot close db connection", "error", err)
	}
}
