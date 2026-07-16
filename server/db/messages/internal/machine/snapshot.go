package machine

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/store"
)

type snapshot struct {
	store  *store.Store
	dumper Dumper
	ctx    context.Context
	ch     <-chan *messages.MessagesSnapshot
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	logger.Debugln("snapshotting messages")

	s := &snapshot{}

	file, err := os.CreateTemp("", "messages_*.bin")
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

		if uint64(n) < size {
			n2, err := s.Read(data[n:])
			if err != nil {
				return err
			}

			logger.Debugw("restore", "n2", n2)
		}

		logger.Debugw("restore", "n", n)

		var snapshot messages.MessagesSnapshot
		if err = proto.Unmarshal(data, &snapshot); err != nil {
			logger.Debugln(zap.ByteString("data", data), zap.Error(err))
			continue
		}

		bytes, err := prototext.Marshal(&snapshot)
		if err != nil {
			logger.Debugln(zap.Error(err))
			continue
		}

		logger.Debugln(zap.String("snapshot", string(bytes)))

		err = f.dumper.Restore(context.TODO(), &snapshot)
		if err != nil {
			logger.Debugln(zap.ByteString("data", data), zap.Error(err))
			continue
		}
	}

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	logger.Debugln("persisting snapshot")

	for snapshot := range s.ch {
		switch v := snapshot.Item.(type) {
		case *messages.MessagesSnapshot_Message:
			logger.Debugw("message snapshot", "id", v.Message.Id)
		case *messages.MessagesSnapshot_Translation:
			logger.Debugw("translation snapshot", "message_id", v.Translation.Id)
		case *messages.MessagesSnapshot_Comment:
			logger.Debugw("comment snapshot", "id", v.Comment.Id)
		default:
			logger.Debugln("unknown snapshot")
			continue
		}

		bytes, err := prototext.Marshal(snapshot)
		if err != nil {
			logger.Debugln(zap.Error(err))
			continue
		}

		logger.Debugln(zap.String("snapshot", string(bytes)))

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
