package saga

import (
	"fmt"
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/sec"
	threadspkg "github.com/bd878/gallery/server/threads/pkg"
	messagespkg "github.com/bd878/gallery/server/messages/pkg"
)

const deleteMessageReplyChannel = "gallery.messages.replies.DeleteMessage"
const DeleteMessageSagaName = "messages.DeleteMessage"

type (
	deleteMessageSaga struct {
		sec.Saga
	}
)

func NewDeleteMessageSaga(nodeName string) sec.Saga {
	saga := deleteMessageSaga{
		Saga: sec.NewSaga(DeleteMessageSagaName, fmt.Sprintf("%s.%s", deleteMessageReplyChannel, nodeName)),
	}

	// 0. -RestoreMessage
	saga.AddStep().
		Compensation(saga.restoreMessage)

	// 1. DeleteThread, -RestoreThread
	saga.AddStep().
		Action(saga.deleteThread).
		Compensation(saga.restoreThread)

	return saga
}

func (s deleteMessageSaga) deleteThread(ctx context.Context, data []byte) am.Command {
	slog.Debug("delete thread saga command")

	return am.NewCommand(threadspkg.DeleteThreadCommand, threadspkg.CommandChannel, data)
}

func (s deleteMessageSaga) restoreMessage(ctx context.Context, data []byte) am.Command {
	slog.Debug("restore message saga command")

	m := &threads.DeleteThread{}
	if err := proto.Unmarshal(data, m); err != nil {
		slog.Error("failed to unmarshal data", slog.String("error", err.Error()))
		return nil
	}

	data2, err := proto.Marshal(&messages.RestoreMessage{
		Id: m.ThreadId,
		UserId: m.UserId,
	})
	if err != nil {
		slog.Error("failed to marshal data", slog.String("error", err.Error()))
		return nil
	}

	return am.NewCommand(messagespkg.RestoreMessageCommand, messagespkg.CommandChannel, data2)
}

func (s deleteMessageSaga) restoreThread(ctx context.Context, data []byte) am.Command {
	slog.Debug("restore thread saga command")

	m := &threads.DeleteThread{}
	if err := proto.Unmarshal(data, m); err != nil {
		slog.Error("failed to unmarshal command", slog.String("error", err.Error()))
		return nil
	}

	data2, err := proto.Marshal(&threads.RestoreThread{
		ThreadId: m.ThreadId,
		UserId: m.UserId,
	})
	if err != nil {
		slog.Error("failed to marshal command", slog.String("error", err.Error()))
		return nil
	}

	return am.NewCommand(threadspkg.RestoreThreadCommand, threadspkg.CommandChannel, data2)
}