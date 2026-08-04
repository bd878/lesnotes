package service

import (
	"context"
	"fmt"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/search"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/db/search/pkg/loadbalance"
	"github.com/bd878/gallery/server/search/pkg/model"
	"github.com/bd878/gallery/server/db/search/pkg/machine"
)

type Config struct {
	RpcAddr  string
}

type Controller struct {
	conf         Config
	client       search.SearchClient
	conn         *grpc.ClientConn
}

func New(conf Config) *Controller {
	controller := &Controller{conf: conf}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			slog.Error("failed to close search conn", slog.String("error", err.Error()))
		}
	}
}

func (s *Controller) setupConnection() (err error) {
	s.Close()

	conn, err := rpc.NewClient(
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

	client := search.NewSearchClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		slog.Debug("search conn failed", slog.String("state", state.String()))
		return true
	}
	return false
}

func (s *Controller) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool, createdAt, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("save search message", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("name", name), slog.String("title", title), slog.String("text", text), slog.Bool("private", private))

	cmd, err := proto.Marshal(&search.AppendMessageCommand{
		Id:          id,
		Text:        text,
		Title:       title,
		Name:        name,
		UserId:      userID,
		Private:     private,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendMessageRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete search message", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.DeleteMessageCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteMessageRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) UpdateMessage(ctx context.Context, id, userID int64, name, title, text *string, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("update search message", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Any("name", name), slog.Any("title", title), slog.Any("text", text))

	cmd, err := proto.Marshal(&search.UpdateMessageCommand{
		Id:         id,
		Text:       text,
		Title:      title,
		Name:       name,
		UserId:     userID,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateMessageRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PublishMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("publish search messages", slog.Any("ids", ids), slog.Int64("user_id", userID), slog.String("updated_at", updatedAt))

	cmd, err := proto.Marshal(&search.PublishMessagesCommand{
		Ids:       ids,
		UserId:    userID,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishMessagesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PrivateMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("private search messages", slog.Any("ids", ids), slog.Int64("user_id", userID), slog.String("updated_at", updatedAt))

	cmd, err := proto.Marshal(&search.PrivateMessagesCommand{
		Ids:         ids,
		UserId:      userID,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateMessagesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool, createdAt, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("save thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.String("name", name), slog.String("description", description), slog.Bool("private", private))

	cmd, err := proto.Marshal(&search.AppendThreadCommand{
		Id:          id,
		Name:        name,
		UserId:      userID,
		ParentId:    parentID,
		Description: description,
		Private:     private,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendThreadRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.DeleteThreadCommand{
		UserId:   userID,
		Id:       id,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteThreadRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) UpdateThread(ctx context.Context, id, userID int64, name, description *string, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("update thread", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Any("name", name), slog.Any("description", description))

	cmd, err := proto.Marshal(&search.UpdateThreadCommand{
		Id:          id,
		Description: description,
		Name:        name,
		UserId:      userID,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateThreadRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) ChangeThreadParent(ctx context.Context, id, userID, parentID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("change thread parent", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Int64("parent_id", parentID))

	cmd, err := proto.Marshal(&search.ChangeThreadParentCommand{
		Id:          id,
		UserId:      userID,
		ParentId:    parentID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.ChangeThreadParentRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PrivateThread(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("private thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.PrivateThreadCommand{
		Id:         id,
		UserId:     userID,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateThreadRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PublishThread(ctx context.Context, id, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("publish thread", slog.Int64("id", id), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.PublishThreadCommand{
		Id:         id,
		UserId:     userID,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishThreadRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64, createdAt, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("save file", slog.Int64("id", id), slog.Int64("user_id", userID), slog.String("name", name), slog.String("description", description), slog.String("mime", mime), slog.Bool("private", private), slog.Int64("size", size))

	cmd, err := proto.Marshal(&search.AppendFileCommand{
		Id:          id,
		UserId:      userID,
		Name:        name,
		Description: description,
		Mime:        mime,
		Private:     private,
		Size:        size,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendFileRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PublishFiles(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("publish files", slog.Any("ids", ids), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.PublishFilesCommand{
		Ids:         ids,
		UserId:      userID,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PublishFilesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) PrivateFiles(ctx context.Context, ids []int64, userID int64, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("private files", slog.Any("ids", ids), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.PrivateFilesCommand{
		Ids:         ids,
		UserId:      userID,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.PrivateFilesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) DeleteFiles(ctx context.Context, ids []int64, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete files", slog.Any("ids", ids), slog.Int64("user_id", userID))

	cmd, err := proto.Marshal(&search.DeleteFilesCommand{
		Ids:         ids,
		UserId:      userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteFilesRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) SaveTranslation(ctx context.Context, userID, messageID int64, lang string, title, text string, createdAt, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("save translation", slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("lang", lang), slog.String("title", title), slog.String("text", text))

	cmd, err := proto.Marshal(&search.AppendTranslationCommand{
		UserId:      userID,
		MessageId:   messageID,
		Lang:        lang,
		Text:        text,
		Title:       title,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("delete translation", slog.Int64("message_id", messageID), slog.String("lang", lang))

	cmd, err := proto.Marshal(&search.DeleteTranslationCommand{
		MessageId:    messageID,
		Lang:         lang,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string, updatedAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("update translation", slog.Int64("message_id", messageID), slog.String("lang", lang), slog.Any("title", title), slog.Any("text", text))

	cmd, err := proto.Marshal(&search.UpdateTranslationCommand{
		MessageId:     messageID,
		Lang:          lang,
		Title:         title,
		Text:          text,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int32) (list []*model.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	slog.Debug("search messages", slog.Int64("user_id", userID), slog.String("substr", substr), slog.Int64("thread_id", threadID), slog.Int("public", int(public)))

	res, err := s.client.SearchMessages(ctx, &search.SearchMessagesRequest{
		Substr:   substr,
		UserId:   userID,
		ThreadId: &threadID,
		Public:   &public,
	})
	if err != nil {
		return nil, err
	}

	list = model.MapMessagesFromProto(model.MessageFromProto, res.List)

	slog.Debug("found messages", slog.Int("count", len(list)))

	return 
}