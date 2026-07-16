package machine

import (
	"io"
	"os"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/billing"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/store"
)

type snapshot struct {
	store         *store.Store
	dumper        Dumper
	ctx           context.Context
	ch            <-chan *billing.BillingSnapshot
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	logger.Debugln("snapshotting billing")

	s := &snapshot{}

	file, err := os.CreateTemp("", "billing_*.bin")
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

		var snapshot billing.BillingSnapshot
		if err = proto.Unmarshal(data, &snapshot); err != nil {
			return err
		}

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
		case *billing.BillingSnapshot_Invoice:
			logger.Debugw("invoice snapshot", "id", v.Invoice.Id)
		case *billing.BillingSnapshot_Payment:
			logger.Debugw("payment snapshot", "id", v.Payment.Id)
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
