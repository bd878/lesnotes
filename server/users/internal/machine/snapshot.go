package machine

import (
	"io"
	"os"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/store"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type snapshot struct {
	store         *store.Store
	usersDumper   UsersDumper
	ctx           context.Context
	ch            <-chan *api.UserSnapshot
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	logger.Debugln("snapshotting users")

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
	logger.Debugln("restoring fsm from snapshot")

	store := store.NewReader(reader)
	defer store.Close()

	for {
		size, err := store.ReadSize()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		data := make([]byte, size)
		n, err := store.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		logger.Debugw("read", "n", n, "data", data)

		var user api.UserSnapshot
		if err = proto.Unmarshal(data, &user); err != nil {
			return err
		}

		logger.Debugw("user", "id", user.Id, "login", user.Login, "salt", user.HashedPassword,
			"created_at", user.CreatedAt, "updated_at", user.UpdatedAt, "metadata", user.Metadata)

		err = f.usersDumper.Restore(context.TODO(), &user)
		if err != nil {
			return err
		}
	}

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	logger.Debugln("persisting snapshot")

	for u := range s.ch {
		if u.Id == model.PublicUserID {
			// restore public user from migration
			continue
		}

		data, err := proto.Marshal(u)
		if err != nil {
			return err
		}

		n, err := s.store.Append(data)
		if err != nil {
			return err
		}

		logger.Debugw("append user to snapshot", "id", u.Id, "n", n)

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

	logger.Debugw("store persisted", "n", n)

	return
}

func (s *snapshot) Release() {
	logger.Debugln("release snapshot")
	if err := s.store.Close(); err != nil {
		logger.Errorw("cannot close store file", "error", err)
	}

	if err := s.usersDumper.Close(); err != nil {
		logger.Errorw("cannot close db connection", "error", err)
	}
}
