package application

import (
	"time"
	"context"
	"bytes"
	"log/slog"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/search"
	"github.com/bd878/gallery/server/db/search/pkg/machine"
)

type MessagesRepository interface {
	SearchMessages(ctx context.Context, userID int64, substr string, public int) (list []*search.SearchMessage, err error)
}

type FilesRepository interface {
}

type TranslationsRepository interface {
	SearchTranslations(ctx context.Context, userID int64, substr string) (list []*search.SearchTranslation, err error)
}

type ThreadsRepository interface {
	SearchThreads(ctx context.Context, parentID, userID int64) (list []*search.SearchThread, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
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

	return m.consensus.Apply(buf.Bytes(), duration)
}

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*search.SearchMessage, err error) {
	slog.Debug("search messages", slog.Int64("user_id", userID), slog.String("substr", substr), slog.Int64("thread_id", threadID), slog.Int("public", public))

	messages, err := m.messagesRepo.SearchMessages(ctx, userID, substr, public)
	if err != nil {
		return nil, err
	}

	slog.Debug("found messages", slog.Int("count", len(messages)))

	if threadID == -1 && threadID == 0 {
		return messages, nil
	}

	// get child threads
	threads, err := m.threadsRepo.SearchThreads(ctx, threadID, userID)
	if err != nil {
		return nil, err
	}

	list = make([]*search.SearchMessage, 0)

	// filter by thread parent id
	for _, msg := range messages {
		for _, thread := range threads {
			if msg.Id == thread.Id {
				list = append(list, msg)
			}
		}
	}

	return
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	slog.Debug("get servers")
	return m.consensus.GetServers(ctx)
}


// TODO: search threads
// TODO: search files
// TODO: search translations