package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
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

func (m *Distributed) apply(ctx context.Context, reqType machine.RequestType, cmd []byte) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), 10*time.Second)
}
func (m *Distributed) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
	// for integration events; though raft will not allow .apply for not a leader, anyway
	// may be we do not need raft for search, when every node receives a message
	m.log.Debugw("save search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text, "private", private)

	cmd, err := proto.Marshal(&machine.AppendMessageCommand{
		Id:      id,
		Text:    text,
		Title:   title,
		Name:    name,
		UserId:  userID,
		Private: private,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendMessageRequest, cmd)

	return
}

func (m *Distributed) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("delete search message", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteMessageCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteMessageRequest, cmd)

	return
}

func (m *Distributed) UpdateMessage(ctx context.Context, id, userID int64, name, title, text string) (err error) {
	m.log.Debugw("update search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text)

	cmd, err := proto.Marshal(&machine.UpdateMessageCommand{
		Id:      id,
		Text:    text,
		Title:   title,
		Name:    name,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateMessageRequest, cmd)

	return
}

func (m *Distributed) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	m.log.Debugw("publish search messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PublishMessagesCommand{
		Ids:     ids,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PublishMessagesRequest, cmd)

	return
}

func (m *Distributed) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	m.log.Debugw("private search messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PrivateMessagesCommand{
		Ids:     ids,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PrivateMessagesRequest, cmd)

	return
}

func (m *Distributed) SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool) (err error) {
	m.log.Debugw("save thread", "id", id, "user_id", userID, "parent_id", parentID, "name", name, "description", description, "private", private)

	cmd, err := proto.Marshal(&machine.AppendThreadCommand{
		Id:          id,
		Name:        name,
		UserId:      userID,
		ParentId:    parentID,
		Description: description,
		Private:     private,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendThreadRequest, cmd)

	return
}

func (m *Distributed) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteThreadCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteThreadRequest, cmd)

	return
}

func (m *Distributed) UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error) {
	m.log.Debugw("update thread", "id", id, "user_id", userID, "name", name, "description", description)

	cmd, err := proto.Marshal(&machine.UpdateThreadCommand{
		Id:          id,
		Description: description,
		Name:        name,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateThreadRequest, cmd)

	return
}

func (m *Distributed) ChangeThreadParent(ctx context.Context, id, userID, parentID int64) (err error) {
	m.log.Debugw("change thread parent", "id", id, "user_id", userID, "parent_id", parentID)

	cmd, err := proto.Marshal(&machine.ChangeThreadParentCommand{
		Id:          id,
		UserId:      userID,
		ParentId:    parentID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.ChangeThreadParentRequest, cmd)

	return
}

func (m *Distributed) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("private thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PrivateThreadCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PrivateThreadRequest, cmd)

	return
}

func (m *Distributed) PublishThread(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("publish thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PublishThreadCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PublishThreadRequest, cmd)

	return
}

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*model.Message, err error) {
	m.log.Debugw("search messages", "user_id", userID, "substr", substr, "thread_id", threadID, "public", public)

	messages, err := m.messagesRepo.SearchMessages(ctx, userID, substr, public)
	if err != nil {
		return nil, err
	}

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

func (m *Distributed) SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64) (err error) {
	m.log.Debugw("save file", "id", id, "user_id", userID, "name", name, "description", description, "mime", mime, "private", private, "size", size)

	cmd, err := proto.Marshal(&machine.AppendFileCommand{
		Id:          id,
		UserId:      userID,
		Name:        name,
		Description: description,
		Mime:        mime,
		Private:     private,
		Size:        size,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendFileRequest, cmd)

	return
}

func (m *Distributed) PublishFile(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("publish file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PublishFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PublishFileRequest, cmd)

	return
}

func (m *Distributed) PrivateFile(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("private file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PrivateFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PrivateFileRequest, cmd)

	return
}

func (m *Distributed) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("delete file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteFileRequest, cmd)

	return
}

func (m *Distributed) SaveTranslation(ctx context.Context, userID, messageID int64, lang string, title, text string) (err error) {
	m.log.Debugw("save translation", "user_id", userID, "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&machine.AppendTranslationCommand{
		UserId:      userID,
		MessageId:   messageID,
		Lang:        lang,
		Text:        text,
		Title:       title,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendTranslationRequest, cmd)

	return
}

func (m *Distributed) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	m.log.Debugw("delete translation", "message_id", messageID, "lang", lang)

	cmd, err := proto.Marshal(&machine.DeleteTranslationCommand{
		MessageId:    messageID,
		Lang:         lang,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteTranslationRequest, cmd)

	return
}

func (m *Distributed) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	m.log.Debugw("update translation", "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&machine.UpdateTranslationCommand{
		MessageId:     messageID,
		Lang:          lang,
		Title:         title,
		Text:          text,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateTranslationRequest, cmd)

	return
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}


// TODO: search threads
// TODO: search files
// TODO: search translations