package application

import (
	"time"
	"context"
	"bytes"
	"log/slog"

	"google.golang.org/protobuf/encoding/prototext"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/search"
	"github.com/bd878/gallery/server/db/search/pkg/machine"
)

type MessagesRepository interface {
	SearchUserMessages(ctx context.Context, userID int64, substr string, public int) (list []*search.SearchMessage, err error)
	SearchBatchMessages(ctx context.Context, userID int64, ids []int64) (list []*search.SearchMessage, err error)
}

type FilesRepository interface {
}

type TranslationsRepository interface {
	SearchTranslations(ctx context.Context, userID int64, substr string) (list []*search.SearchTranslation, err error)
}

type ThreadsRepository interface {
	ReadThread(ctx context.Context, name string) (thread *search.SearchThread, err error)
	SearchThreads(ctx context.Context, parentID, userID int64) (list []*search.SearchThread, err error)
	SearchThreadIDs(ctx context.Context, parentID, userID int64) (ids []int64, err error)
}

type Consensus interface {
	Apply(ctx context.Context, cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	messagesRepo      MessagesRepository
	threadsRepo       ThreadsRepository
	filesRepo         FilesRepository
	translationsRepo  TranslationsRepository
}

func New(consensus Consensus, messagesRepo MessagesRepository,
	filesRepo FilesRepository, threadsRepo ThreadsRepository, translationsRepo TranslationsRepository) *Distributed {
	return &Distributed{
		consensus:        consensus,
		messagesRepo:     messagesRepo,
		threadsRepo:      threadsRepo,
		filesRepo:        filesRepo,
		translationsRepo: translationsRepo,
	}
}

func (m *Distributed) Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	timedCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	return m.consensus.Apply(timedCtx, buf.Bytes(), duration)
}

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, name string, public int) (list []*search.SearchMessage, err error) {
	slog.Debug("search messages",
		slog.Int64("user_id", userID),
		slog.String("substr", substr),
		slog.Int64("thread_id", threadID),
		slog.Int("public", public),
		slog.String("name", name),
	)

	if threadID != 0 || name != "" {
		if threadID == 0 && name != "" {
			thread, err := m.threadsRepo.ReadThread(ctx, name)
			if err != nil {
				return nil, err
			}
			threadID = thread.Id
		}

		ids, err := m.threadsRepo.SearchThreadIDs(ctx, threadID, userID)
		if err != nil {
			return nil, err
		}

		return m.messagesRepo.SearchBatchMessages(ctx, userID, ids)
	}

	return m.messagesRepo.SearchUserMessages(ctx, userID, substr, public)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	slog.Debug("get servers")
	servers, err := m.consensus.GetServers(ctx)
	if err != nil {
		return nil, err
	}
	out, err := prototext.Marshal(&api.Servers{Servers: servers})
	if err != nil {
		return nil, err
	}

	slog.Debug("servers", string(out))
	return servers, nil
}


// TODO: search threads
// TODO: search files
// TODO: search translations