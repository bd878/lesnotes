package service

import (
	"context"
	"fmt"
	"sync"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/pkg/model"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type MessagesConfig struct {
	RpcAddr string
}

type ThreadsGateway interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32) (list []*threads.Thread, isLastPage bool, err error)
	ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, privateMessage *bool) (list []*threads.Thread, isLastPage bool, err error)
	// TODO: fix params order
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error)
	CreateThread(ctx context.Context, id, userID, parentID int64, name string, private bool) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	UpdateThread(ctx context.Context, id, userID int64) (err error)
	ResolvePath(ctx context.Context, userID, id int64) (path []*api.PathStep, err error)
	ReadThread(ctx context.Context, userID, id int64, name string) (thread *threads.Thread, err error)
	ReadParent(ctx context.Context, userID, id int64) (thread *threads.Thread, err error)
}

type MessagesController struct {
	conf    MessagesConfig
	client  api.MessagesClient
	conn    *grpc.ClientConn
	threads ThreadsGateway
}

func NewMessagesController(conf MessagesConfig, threads ThreadsGateway) *MessagesController {
	controller := &MessagesController{conf: conf, threads: threads}

	controller.setupConnection()

	return controller
}

func (s *MessagesController) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *MessagesController) setupConnection() (err error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := api.NewMessagesClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *MessagesController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugw("connection failed", "state", state.String())
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

	logger.Debugw("save message", "id", id, "text", text, "title", title, "file_ids", fileIDs, "thread_id", threadID, "user_id", userID, "private", private, "name", name)

	_, err = s.client.SaveMessage(ctx, &api.SaveMessageRequest{
		Id:      id,
		Text:    text,
		Title:   title,
		FileIds: fileIDs,
		UserId:  userID,
		Private: private,
		Name:    name,
	})
	if err != nil {
		return
	}

	// TODO: what if we fall here? Message is created, but thread is not yet.
	// In readThreadMessages we request by batch. We arrive in situation, where
	// len(threads) != len(messages).

	err = s.threads.CreateThread(ctx, id, userID, threadID, name, private)
	if err != nil {
		if status, ok := status.FromError(err); ok {
			logger.Debugw("thread failed", "message", status.Message())
		} else {
			logger.Debugw("non-rpc error", "error", err)
		}
		return
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

	return
}

func (s *MessagesController) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete messages", "ids", ids, "user_id", userID)

	// TODO: DeleteThreads
	for _, id := range ids {
		err = s.threads.DeleteThread(ctx, id, userID)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				logger.Debugw("thread failed", "message", status.Message())
			} else {
				logger.Debugw("non-rpc error", "error", err)
			}
			return
		}
	}

	_, err = s.client.DeleteMessages(ctx, &api.DeleteMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *MessagesController) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("publish messages", "ids", ids, "user_id", userID)

	_, err = s.client.PublishMessages(ctx, &api.PublishMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *MessagesController) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("private messages", "ids", ids, "user_id", userID)

	_, err = s.client.PrivateMessages(ctx, &api.PrivateMessagesRequest{
		Ids:    ids,
		UserId: userID,
	})

	return
}

func (s *MessagesController) UpdateMessage(ctx context.Context, id int64, text, title, name *string, fileIDs []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "user_id", userID)

	_, err = s.client.UpdateMessage(ctx, &api.UpdateMessageRequest{
		Id:      id,
		UserId:  userID,
		FileIds: fileIDs,
		Text:    text,
		Title:   title,
		Name:    name,
	})

	return
}

// Get messages in order
func (s *MessagesController) ReadThreadMessages(ctx context.Context, userID, threadID int64, threadName string,
	limit, offset int32, ascending bool, privateMessage *bool) (list *model.MessagesList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread messages", "user_id", userID, "thread_id", threadID, "thread_name",
		threadName, "limit", limit, "offset", offset, "ascending", ascending, "private_message", privateMessage)

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

	logger.Debugw("read thread messages", "count total", total)

	threadsList, isLastPage, err := s.threads.ListMessages(ctx, userID, threadID, limit, offset, privateMessage)
	if err != nil {
		logger.Debugw("failed to list messages", "error", err)
		return nil, err
	}

	logger.Debugw("read thread messages", "threads_list", threadsList)

	ids := make([]int64, 0)
	for _, thread := range threadsList {
		ids = append(ids, thread.ID)
	}

	messages, err := s.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		logger.Debugw("failed to read batch messages", "error", err)
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
func (s *MessagesController) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read batch messages", "user_id", userID, "ids", ids)

	res, err := s.client.ReadBatchMessages(ctx, &api.ReadBatchMessagesRequest{
		UserId: userID,
		Ids:    ids,
	})
	if err != nil {
		return nil, err
	}

	messages = model.MapMessagesFromProto(model.MessageFromProto, res.Messages)

	return
}

// read all messages not concerning thread
func (s *MessagesController) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (list *model.MessagesList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)

	res, err := s.client.ReadMessages(ctx, &api.ReadMessagesRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
		Asc:    ascending,
	})
	if err != nil {
		return nil, err
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

	logger.Debugw("read message", "id", id, "name", name, "user_ids", userIDs)

	res, err := s.client.ReadMessage(ctx, &api.ReadMessageRequest{
		Id:      id,
		UserIds: userIDs,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}

	// TODO: add threads count

	message = model.MessageFromProto(res)

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

	parentMessage, err := s.client.ReadMessage(ctx, &api.ReadMessageRequest{
		Id: parent.ID,
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	message.ParentMessage = &model.Identity{
		ID: parentMessage.Id,
		Title: parentMessage.Title,
		Name: parentMessage.Name,
		Private: parentMessage.Private,
	}

	return
}

func (s *MessagesController) ReadPath(ctx context.Context, userID, id int64, name string) (path []*model.Message, parentID int64, err error) {
	if s.isConnFailed() {
		if err := s.setupConnection(); err != nil {
			return nil, 0, err
		}
	}

	logger.Debugw("read path", "user_id", userID, "id", id, "name", name)

	if name != "" && id == 0 {
		message, err := s.ReadMessage(ctx, id, name, []int64{userID})
		if err != nil {
			return nil, 0, err
		}

		id = message.ID
		logger.Debugw("read path", "message", message)
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

	go func(ctx context.Context, ch chan *model.IDLimitOffset) {

		ps.mu.RLock()
		defer ps.mu.RUnlock()
		defer close(ch)

		for _, pair := range ps.dict {
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

	logger.Debugw("read tree", "user_id", userID, "highlight_id", highlightID, "highlight_name", highlightName, "message_id",
		rootID, "name", rootName, "limit", limit, "offset", offset, "private_message", privateMessage, "pairs", pairs)

	map1 := NewLimitOffsetMap(pairs)

	logger.Debugw("read_tree", "map1", map1.String())

	highlightPath, _, err := s.ReadPath(ctx, userID, highlightID, highlightName)
	if err != nil {
		return nil, err
	}

	logger.Debugw("read_tree", "highlight path", highlightPath)

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