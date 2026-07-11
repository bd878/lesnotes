package application

import (
	"time"
	"context"
	"bytes"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/search/internal/machine"
	"github.com/bd878/gallery/server/search/pkg/model"
)

type MessagesRepository interface {
	SearchMessages(ctx context.Context, userID int64, substr string, public int) (list []*model.Message, err error)
}

type FilesRepository interface {
}

type TranslationsRepository interface {
	SearchTranslations(ctx context.Context, userID int64, substr string) (list []*model.Translation, err error)
}

type ThreadsRepository interface {
	SearchThreads(ctx context.Context, parentID, userID int64) (list []*model.Thread, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	log               *logger.Logger
	messagesRepo      MessagesRepository
	threadsRepo       ThreadsRepository
	filesRepo         FilesRepository
	translationsRepo  TranslationsRepository
}

func New(consensus Consensus, messagesRepo MessagesRepository,
	filesRepo FilesRepository, threadsRepo ThreadsRepository, translationsRepo TranslationsRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:              log,
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

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*model.Message, err error) {
	m.log.Debugw("search messages", "user_id", userID, "substr", substr, "thread_id", threadID, "public", public)

	messages, err := m.messagesRepo.SearchMessages(ctx, userID, substr, public)
	if err != nil {
		return nil, err
	}

	m.log.Debugw("found messages", "count", len(messages))

	if threadID == -1 && threadID == 0 {
		return messages, nil
	}

	// get child threads
	threads, err := m.threadsRepo.SearchThreads(ctx, threadID, userID)
	if err != nil {
		return nil, err
	}

	list = make([]*model.Message, 0)

	// filter by thread parent id
	for _, msg := range messages {
		for _, thread := range threads {
			if msg.ID == thread.ID {
				list = append(list, msg)
			}
		}
	}

	return
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}


// TODO: search threads
// TODO: search files
// TODO: search translations