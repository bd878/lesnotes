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

const createMessageReplyChannel = messageReplyChannel + ".CreateMessage"
const CreateMessageSagaName = "messages.CreateMessage"

type (
	createMessageSaga struct {
		sec.Saga
	}
)

func NewCreateMessageSaga(nodeName string) sec.Saga {
	saga := createMessageSaga{
		Saga: sec.NewSaga(CreateMessageSagaName, fmt.Sprintf("%s.%s", createMessageReplyChannel, nodeName)),
	}

	// 0. -RemoveMessage
	saga.AddStep().
		Compensation(saga.removeMessage)

	// 1. CreateThread, -DeleteThread
	saga.AddStep().
		Action(saga.createThread).
		Compensation(saga.deleteThread)

	return saga
}

func (s createMessageSaga) createThread(ctx context.Context, data []byte) am.Command {
	slog.Debug("create thread saga command")

	return am.NewCommand(threadspkg.CreateThreadCommand, threadspkg.CommandChannel, data)
}

func (s createMessageSaga) removeMessage(ctx context.Context, data []byte) am.Command {
	slog.Debug("remove message saga command")

	m := &threads.CreateThread{}
	if err := proto.Unmarshal(data, m); err != nil {
		slog.Error("failed to unmarshal data", slog.String("error", err.Error()))
		return nil
	}

	data2, err := proto.Marshal(&messages.DeleteMessage{
		Id: m.ThreadId,
		UserId: m.UserId,
	})
	if err != nil {
		slog.Error("failed to marshal data", slog.String("error", err.Error()))
		return nil
	}

	return am.NewCommand(messagespkg.DeleteMessageCommand, messagespkg.CommandChannel, data2)
}

func (s createMessageSaga) deleteThread(ctx context.Context, data []byte) am.Command {
	slog.Debug("delete thread saga command")

	m := &threads.CreateThread{}
	if err := proto.Unmarshal(data, m); err != nil {
		slog.Error("failed to unmarshal command", slog.String("error", err.Error()))
		return nil
	}

	data2, err := proto.Marshal(&threads.DeleteThread{
		ThreadId: m.ThreadId,
		UserId: m.UserId,
	})
	if err != nil {
		slog.Error("failed to marshal command", slog.String("error", err.Error()))
		return nil
	}

	return am.NewCommand(threadspkg.DeleteThreadCommand, threadspkg.CommandChannel, data2)
}