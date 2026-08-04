package machine

import (
	"io"
	"os"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/users"
	"log/slog"
	"github.com/bd878/gallery/server/internal/store"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type snapshot struct {
	store         *store.Store
	usersDumper   UsersDumper
	ctx           context.Context
	ch            <-chan *users.UsersSnapshot
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	slog.Debug("snapshotting users")

	s := &snapshot{}

	file, err := os.CreateTemp("", "users_*.bin")
	if err != nil {
		return nil, err
	}

	s.store, err = store.NewStore(file)
	if err != nil {
		return nil, err
	}

	s.usersDumper = f.usersDumper
	s.ctx = context.TODO()
	s.ch, err = f.usersDumper.Open(s.ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (f *Machine) Restore(reader io.ReadCloser) (err error) {
	slog.Debug("restoring fsm from snapshot")

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

		data := make([]byte, size)
		n, err := s.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if uint64(n) < size {
			_, err := s.Read(data[n:])
			if err != nil {
				return err
			}
		}

		var snapshot users.UsersSnapshot
		if err = proto.Unmarshal(data, &snapshot); err != nil {
			return err
		}

		err = f.usersDumper.Restore(context.TODO(), &snapshot)
		if err != nil {
			return err
		}
	}

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	slog.Debug("persisting snapshot")

	for snapshot := range s.ch {
		switch snapshot.Item.(type) {
		case *users.UsersSnapshot_User:
		case *users.UsersSnapshot_Premium:
		default:
			slog.Error("unknown snapshot")
			continue
		}

		if u, ok := snapshot.Item.(*users.UsersSnapshot_User); ok {
			if u.User.Id == model.PublicUserID {
				// restore public user from migration
				continue
			}
		}

		data, err := proto.Marshal(snapshot)
		if err != nil {
			return err
		}

		_, err = s.store.Append(data)
		if err != nil {
			return err
		}

		select {
		case <-s.ctx.Done():
			return context.Cause(s.ctx)
		default:
		}
	}

	err = s.store.Seek()
	if err != nil {
		return
	}

	n, err := io.Copy(sink, s.store.File)
	if err != nil {
		return err
	}

	slog.Debug("store persisted", slog.Int64("n", n))

	return
}

func (s *snapshot) Release() {
	slog.Debug("release snapshot")
	if err := s.store.Close(); err != nil {
		slog.Error("cannot close store file", slog.String("error", err.Error()))
	}

	if err := s.usersDumper.Close(); err != nil {
		slog.Error("cannot close db connection", slog.String("error", err.Error()))
	}
}
