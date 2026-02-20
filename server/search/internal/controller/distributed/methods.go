package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/search/pkg/model"
)

func (m *Distributed) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	timeout := 10*time.Second
	/* fsm.Apply() */
	future := m.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res = future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}

	return
}

func (m *Distributed) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
	// for integration events; though raft will not allow .apply for not a leader, anyway
	// may be we do not need raft for search, when every node receives a message
	if !m.isLeader() {
		return
	}

	logger.Debugw("save search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text, "private", private)

	cmd, err := proto.Marshal(&AppendMessageCommand{
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

	_, err = m.apply(ctx, AppendMessageRequest, cmd)

	return
}

func (m *Distributed) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return
	}

	logger.Debugw("delete search message", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteMessageCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteMessageRequest, cmd)

	return
}

func (m *Distributed) UpdateMessage(ctx context.Context, id, userID int64, name, title, text string) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("update search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text)

	cmd, err := proto.Marshal(&UpdateMessageCommand{
		Id:      id,
		Text:    text,
		Title:   title,
		Name:    name,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, UpdateMessageRequest, cmd)

	return
}

func (m *Distributed) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("publish search messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&PublishMessagesCommand{
		Ids:     ids,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishMessagesRequest, cmd)

	return
}

func (m *Distributed) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("private search messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateMessagesCommand{
		Ids:     ids,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateMessagesRequest, cmd)

	return
}

func (m *Distributed) SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("save thread", "id", id, "user_id", userID, "parent_id", parentID, "name", name, "description", description, "private", private)

	cmd, err := proto.Marshal(&AppendThreadCommand{
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

	_, err = m.apply(ctx, AppendThreadRequest, cmd)

	return
}

func (m *Distributed) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("delete thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteThreadCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteThreadRequest, cmd)

	return
}

func (m *Distributed) UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("update thread", "id", id, "user_id", userID, "name", name, "description", description)

	cmd, err := proto.Marshal(&UpdateThreadCommand{
		Id:          id,
		Description: description,
		Name:        name,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, UpdateThreadRequest, cmd)

	return
}

func (m *Distributed) ChangeThreadParent(ctx context.Context, id, userID, parentID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("change thread parent", "id", id, "user_id", userID, "parent_id", parentID)

	cmd, err := proto.Marshal(&ChangeThreadParentCommand{
		Id:          id,
		UserId:      userID,
		ParentId:    parentID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, ChangeThreadParentRequest, cmd)

	return
}

func (m *Distributed) PrivateThread(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("private thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateThreadCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateThreadRequest, cmd)

	return
}

func (m *Distributed) PublishThread(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("publish thread", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PublishThreadCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishThreadRequest, cmd)

	return
}

func (m *Distributed) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*model.Message, err error) {
	logger.Debugw("search messages", "user_id", userID, "substr", substr, "thread_id", threadID, "public", public)

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
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("save file", "id", id, "user_id", userID, "name", name, "description", description, "mime", mime, "private", private, "size", size)

	cmd, err := proto.Marshal(&AppendFileCommand{
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

	_, err = m.apply(ctx, AppendFileRequest, cmd)

	return
}

func (m *Distributed) PublishFile(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("publish file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PublishFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PublishFileRequest, cmd)

	return
}

func (m *Distributed) PrivateFile(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("private file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&PrivateFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, PrivateFileRequest, cmd)

	return
}

func (m *Distributed) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("delete file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&DeleteFileCommand{
		Id:          id,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteFileRequest, cmd)

	return
}

func (m *Distributed) SaveTranslation(ctx context.Context, userID, messageID int64, lang string, title, text string) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("save translation", "user_id", userID, "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&AppendTranslationCommand{
		UserId:      userID,
		MessageId:   messageID,
		Lang:        lang,
		Text:        text,
		Title:       title,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, AppendTranslationRequest, cmd)

	return
}

func (m *Distributed) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("delete translation", "message_id", messageID, "lang", lang)

	cmd, err := proto.Marshal(&DeleteTranslationCommand{
		MessageId:    messageID,
		Lang:         lang,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteTranslationRequest, cmd)

	return
}

func (m *Distributed) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	if !m.isLeader() {
		return nil
	}

	logger.Debugw("update translation", "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&UpdateTranslationCommand{
		MessageId:     messageID,
		Lang:          lang,
		Title:         title,
		Text:          text,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, UpdateTranslationRequest, cmd)

	return
}

// TODO: search threads
// TODO: search files
// TODO: search translations