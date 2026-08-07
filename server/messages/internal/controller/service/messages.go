package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/files"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
	"github.com/bd878/gallery/server/db/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/pkg/model"
	threadsmodel "github.com/bd878/gallery/server/threads/pkg/model"
)

type MessagesConfig struct {
	RpcAddr string
}

type ThreadsGateway interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32) (list []*threadsmodel.Thread, isLastPage bool, err error)
	ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, privateMessage *bool) (list []*threadsmodel.Thread, isLastPage bool, err error)
	// TODO: fix params order
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error)
	ResolvePath(ctx context.Context, userID, id int64) (path []*threads.PathStep, err error)
	ReadThread(ctx context.Context, userID, id int64, name string) (thread *threadsmodel.Thread, err error)
	ReadParent(ctx context.Context, userID, id int64) (thread *threadsmodel.Thread, err error)
}

type FilesGateway interface {
	ReadMessageFiles(ctx context.Context, id int64, userIDs []int64) (list []*files.File, err error)
}

type MessagesController struct {
	conf    MessagesConfig
	client  messages.MessagesClient
	conn    *grpc.ClientConn
	threads ThreadsGateway
	filesGateway FilesGateway
	messagesSaved prometheus.Counter
	publisher    ddd.EventPublisher[ddd.Event]
}

func NewMessagesController(conf MessagesConfig, publisher ddd.EventPublisher[ddd.Event], filesGateway FilesGateway,
	threads ThreadsGateway, messagesSaved prometheus.Counter) *MessagesController {
	controller := &MessagesController{
		conf: conf,
		filesGateway: filesGateway,
		threads: threads,
		publisher: publisher,
		messagesSaved: messagesSaved,
	}

	controller.setupConnection()

	return controller
}

func (s *MessagesController) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			slog.Error("failed to close messages conn", slog.String("error", err.Error()))
		}
	}
}

func (s *MessagesController) setupConnection() (err error) {
	s.Close()

	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
	)
	if err != nil {
		return err
	}

	client := messages.NewMessagesClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *MessagesController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		slog.Debug("messages conn failed", slog.String("state", state.String()))
		return true
	}
	return false
}

func (s *MessagesController) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("save message", slog.Int64("id", id), slog.String("text", text), slog.String("title", title), slog.String("file_ids", fmt.Sprintf("%v", fileIDs)), slog.Int64("thread_id", threadID), slog.Int64("user_id", userID), slog.Bool("private", private), slog.String("name", name))

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&messages.AppendCommand{
		Id:        id,
		Text:      text,
		Title:     title,
		FileIds:   fileIDs,
		UserId:    userID,
		Private:   private,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	event, err := domain.CreateMessage(id, threadID, text, title, fileIDs, userID, private, name, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return nil, err
	}

	message = &model.Message{
		ID:      id,
		Text:    text,
		Title:   title,
		Name:    name,
		FileIDs: fileIDs,
		UserID:  userID,
		Private: private,
	}

	s.messagesSaved.Inc()

	return
}

func (s *MessagesController) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete user messages", slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&messages.DeleteUserMessagesCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteUserMessagesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *MessagesController) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete messages", slog.String("ids", fmt.Sprintf("%v", ids)), slog.Int64("user_id", userID))

	for _, id := range ids {

		event, err := domain.DeleteMessage(id, userID)
		if err != nil {
			return err
		}

		cmd, err := proto.Marshal(&messages.DeleteCommand{
			Id:     id,
			UserId: userID,
		})
		if err != nil {
			return err
		}

		_, err = s.client.Apply(ctx, &api.Command{
			ReqType: int32(machine.DeleteRequest),
			Cmd: cmd,
			Duration: "10s",
		})
		if err != nil {
			return err
		}

		s.publisher.Publish(ctx, event)

	}

	return
}

func (s *MessagesController) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("publish messages", slog.String("ids", fmt.Sprintf("%v", ids)), slog.Int64("user_id", userID))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&messages.PublishCommand{
		Ids:       ids,
		UserId:    userID,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.PublishMessages(userID, ids, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *MessagesController) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("private messages", slog.String("ids", fmt.Sprintf("%v", ids)), slog.Int64("user_id", userID))

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&messages.PrivateCommand{
		Ids:       ids,
		UserId:    userID,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.PrivateMessages(userID, ids, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *MessagesController) UpdateMessage(ctx context.Context, id int64, text, title, name *string, fileIDs []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logValues := []any{slog.Int64("id", id), slog.Int64("user_id", userID)}
	if text != nil {
		logValues = append(logValues, slog.String("text", *text))
	}
	if title != nil {
		logValues = append(logValues, slog.String("title", *title))
	}
	if name != nil {
		logValues = append(logValues, slog.String("name", *name))
	}
	logValues = append(logValues, slog.String("file_ids", fmt.Sprintf("%v", fileIDs)))
	slog.Debug("update message", logValues...)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateMessage(id, text, title, fileIDs, userID, name, updatedAt)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&messages.UpdateCommand{
		Id:        id,
		UserId:    userID,
		FileIds:   fileIDs,
		Text:      text,
		Name:      name,
		Title:     title,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	return s.publisher.Publish(ctx, event)
}

// Get messages in order
func (s *MessagesController) ReadThreadMessages(ctx context.Context, userID, threadID int64, threadName string,
	limit, offset int32, ascending bool, privateMessage *bool) (list *model.MessagesList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logValues := []any{slog.Int64("user_id", userID), slog.Int64("thread_id", threadID), slog.String("thread_name", threadName), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("ascending", ascending)}
	if privateMessage != nil {
		logValues = append(logValues, slog.Bool("private_message", *privateMessage))
	}
	slog.Debug("read thread messages", logValues...)

	// read thread that is not root
	if threadName != "" || threadID != 0 {
		thread, err := s.threads.ReadThread(ctx, userID, threadID, threadName)
		if err != nil {
			return nil, err
		}

		threadID = thread.ID
		userID = thread.UserID
		threadName = thread.Name
	}

	total, err := s.threads.CountMessages(ctx, threadID, userID, privateMessage)
	if err != nil {
		return nil, err
	}

	slog.Debug("read thread messages", slog.Any("count total", total))

	threadsList, isLastPage, err := s.threads.ListMessages(ctx, userID, threadID, limit, offset, privateMessage)
		if err != nil {
			slog.Debug("failed to list messages", slog.Any("error", err))
			return nil, err
		}

	slog.Debug("read thread messages", slog.String("threads_list", fmt.Sprintf("%v", threadsList)))

	ids := make([]int64, 0)
	for _, thread := range threadsList {
		ids = append(ids, thread.ID)
	}

	messages, err := s.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		slog.Debug("failed to read batch messages", slog.Any("error", err))
		return nil, err
	}

	for i, message := range messages {
		message.Count = threadsList[i].Count

		thread, err := s.threads.ReadThread(ctx, message.UserID, message.ID, "")
		if err != nil {
			return nil, err
		}

		// TODO: move to message mapper, this mapping is duplicated here
		message.Thread = &model.Identity{
			ID: thread.ID,
			Title: thread.Title,
			Private: thread.Private,
			Name: thread.Name,
		}
	}

	list = &model.MessagesList{
		Messages:    messages,
		Name:        threadName,
		IsLastPage:  isLastPage,
		IsFirstPage: offset == 0,
		Total:       total,
		Count:       int32(len(ids)),
		Offset:      offset,
	}

	return
}

// Read messages by given ids
func (s *MessagesController) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (list []*model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("read batch messages", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

	res, err := s.client.ReadBatchMessages(ctx, &messages.ReadBatchMessagesRequest{
		UserId: userID,
		Ids:    ids,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			// here it fails : context canceled while waiting for connections to become ready
			slog.Error("read batch messages", slog.String("rpc error", status.Message()))
		} else {
			slog.Error("read batch messages", slog.Any("non-rpc error", err))
		}
		return nil, err
	}

	for _, message := range res.Messages {
		message.Files, err = s.filesGateway.ReadMessageFiles(ctx, message.Id, []int64{userID, message.UserId})
		if err != nil {
			return
		}
	}

	list = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)

	return
}

// read all messages not concerning thread
func (s *MessagesController) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (list *model.MessagesList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("read messages", slog.Int64("user_id", userID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("ascending", ascending))

	res, err := s.client.ReadMessages(ctx, &messages.ReadMessagesRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
		Asc:    ascending,
	})
	if err != nil {
		return nil, err
	}

	for _, message := range res.Messages {
		message.Files, err = s.filesGateway.ReadMessageFiles(ctx, message.Id, []int64{userID, message.UserId})
		if err != nil {
			return
		}
	}

	list = &model.MessagesList{
		Messages:    model.MapMessagesFromProto(model.MessageFromProto, res.Messages),
		IsLastPage:  res.IsLastPage,
		IsFirstPage: offset == 0,
		Offset:      offset,
		// TODO: Total:        total.Count threads.CountThreads,
		Count: int32(len(res.Messages)),
	}

	return
}

func (s *MessagesController) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *model.Message, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	slog.Debug("read message", slog.Int64("id", id), slog.String("name", name), slog.String("user_ids", fmt.Sprintf("%v", userIDs)))

	res, err := s.client.ReadMessage(ctx, &messages.ReadMessageRequest{
		Id:      id,
		UserIds: userIDs,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}

	res.Message.Files, err = s.filesGateway.ReadMessageFiles(ctx, res.Message.Id, append(userIDs, res.Message.UserId))
	if err != nil {
		return
	}

	// TODO: add threads count

	message = model.MessageFromProto(res.Message)

	parent, err := s.threads.ReadParent(ctx, message.UserID, message.ID)
	if err != nil {
		return nil, err
	}

	message.ParentThread = &model.Identity{
		ID:  parent.ID,
		Title: parent.Title,
		Name: parent.Name,
		Private: parent.Private,
	}

	parentMessage, err := s.client.ReadMessage(ctx, &messages.ReadMessageRequest{
		Id: parent.ID,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	if parentMessage.IsRoot {
		message.ParentMessage = &model.Identity{
			ID: 0,
			Title: "",
			Name: "",
			Private: true,
		}
	} else {
		message.ParentMessage = &model.Identity{
			ID: parentMessage.Message.Id,
			Title: parentMessage.Message.Title,
			Name: parentMessage.Message.Name,
			Private: parentMessage.Message.Private,
		}
	}

	return
}

func (s *MessagesController) ReadPath(ctx context.Context, userID, id int64, name string) (path []*model.Message, parentID int64, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, 0, err
		}
	}

	slog.Debug("read path", slog.Int64("user_id", userID), slog.Int64("id", id), slog.String("name", name))

	if name != "" && id == 0 {
		message, err := s.ReadMessage(ctx, id, name, []int64{userID})
		if err != nil {
			return nil, 0, err
		}

		id = message.ID
			slog.Debug("read path", slog.String("message", fmt.Sprintf("%v", message)))
	}

	threads, err := s.threads.ResolvePath(ctx, userID, id)
	if err != nil {
		return nil, 0, err
	}

	// TODO: log error, falls if error 0
	thread, _ := s.threads.ReadThread(ctx, userID, id, "")
	if thread != nil {
		parentID = thread.ParentID
	}

	ids := make([]int64, len(threads))
	for _, step := range threads {
		ids = append(ids, step.Id)
	}

	path, err = s.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		return nil, 0, err
	}

	for i, message := range path {
		message.Thread = &model.Identity{
			ID: threads[i].Id,
			Title: threads[i].Title,
			Name: threads[i].Name,
			Private: threads[i].Private,
		}
	}

	return
}

type limitOffsetMap struct {
	dict map[int64]*model.IDLimitOffset
	mu sync.RWMutex
}

func NewLimitOffsetMap(pairs []*model.IDLimitOffset) (s *limitOffsetMap) {
	s = &limitOffsetMap{
		dict: make(map[int64]*model.IDLimitOffset),
	}

	for _, pair := range pairs {
		s.Add(pair)
	}

	return
}

func (ps *limitOffsetMap) Add(pair *model.IDLimitOffset) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.dict[pair.ID]
	if !ok {
		ps.dict[pair.ID] = pair
	}
}

func (ps *limitOffsetMap) Has(id int64) (ok bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok = ps.dict[id]

	return
}

func (ps *limitOffsetMap) Get(id int64) (result *model.IDLimitOffset, ok bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	result, ok = ps.dict[id]

	return
}

func (ps *limitOffsetMap) Iterator(ctx context.Context) (result chan *model.IDLimitOffset) {
	result = make(chan *model.IDLimitOffset, 0)

	ps.mu.RLock()
	pairs := make([]*model.IDLimitOffset, 0, len(ps.dict))
	for _, pair := range ps.dict {
		pairs = append(pairs, pair)
	}
	ps.mu.RUnlock()

	go func(ctx context.Context, ch chan *model.IDLimitOffset) {
		defer close(ch)
		for _, pair := range pairs {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- pair
			}
		}
	}(ctx, result)

	return
}

func (ps *limitOffsetMap) Union(ctx context.Context, map2 *limitOffsetMap) {
	for pair := range map2.Iterator(ctx) {
		ps.Add(pair)
	}
}

func (ps *limitOffsetMap) String() (result string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var b strings.Builder
	for _, pair := range ps.dict {
		fmt.Fprintf(&b, "{%d, %d, %d}", pair.ID, pair.Limit, pair.Offset)
	}

	return b.String()
}

type idMap struct {
	dict map[int64]struct{}
	mu sync.RWMutex
}

func NewIDMap(list []*model.Message) (m *idMap) {
	m = &idMap{
		dict: make(map[int64]struct{}),
	}

	for _, v := range list {
		m.Add(v)
	}

	return m
}

func (m *idMap) Add(value *model.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.dict[value.ID]
	if !ok {
		m.dict[value.ID] = struct{}{}
	}
}

func (m *idMap) Has(id int64) (ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok = m.dict[id]

	return
}

func (m *idMap) String() (result string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var b strings.Builder
	for k, _ := range m.dict {
		fmt.Fprintf(&b, "{%d}", k)
	}

	return b.String()
}

const limitPathMessages = 10 /* any other number ? */

func (s *MessagesController) ReadTree(ctx context.Context, userID, highlightID int64, highlightName string,
	rootID int64, rootName string, limit, offset int32, privateMessage *bool, pairs []*model.IDLimitOffset) (list *model.MessagesList, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, err
		}
	}

	logValues := []any{slog.Int64("user_id", userID), slog.Int64("highlight_id", highlightID), slog.String("highlight_name", highlightName), slog.Int64("message_id", rootID), slog.String("name", rootName), slog.Int("limit", int(limit)), slog.Int("offset", int(offset))}
	if privateMessage != nil {
		logValues = append(logValues, slog.Bool("private_message", *privateMessage))
	}
	logValues = append(logValues, slog.String("pairs", fmt.Sprintf("%v", pairs)))
	slog.Debug("read tree", logValues...)

	map1 := NewLimitOffsetMap(pairs)

	slog.Debug("read_tree", slog.String("map1", map1.String()))

	highlightPath, _, err := s.ReadPath(ctx, userID, highlightID, highlightName)
	if err != nil {
		return nil, err
	}

	slog.Debug("read_tree", slog.String("highlight_path", fmt.Sprintf("%v", highlightPath)))

	highlightMap := NewIDMap(highlightPath)

	list, err = s.ReadThreadMessages(ctx, userID, rootID, rootName, limit, offset, true, privateMessage)
	if err != nil {
		return nil, err
	}

	err = s.resolveTree(ctx, highlightMap, list, map1, privateMessage)

	return
}

func (s *MessagesController) resolveTree(ctx context.Context, highlightMap *idMap, list *model.MessagesList, map1 *limitOffsetMap, privateMessage *bool) (err error) {
	for _, message := range list.Messages {
		if highlightMap.Has(message.ID) {
			message.Highlight = true
		}

		limitOffset, ok := map1.Get(message.ID)
		if !ok {
			// exit, when there are no more open threads
			continue
		}

		message.Messages, err = s.ReadThreadMessages(ctx, message.UserID, message.ID, "", limitOffset.Limit, limitOffset.Offset, true, privateMessage)
		if err != nil {
			return
		}

		thread, err := s.threads.ReadThread(ctx, message.UserID, message.ID, "")
		if err != nil {
			return err
		}

		message.Thread = &model.Identity{
			ID: thread.ID,
			Title: thread.Title,
			Private: thread.Private,
			Name: thread.Name,
		}

		err = s.resolveTree(ctx, highlightMap, message.Messages, map1, privateMessage)
		if err != nil {
			return err
		}
	}

	return
}