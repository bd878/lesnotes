package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/logger"
)

type Repository interface {
	Create(ctx context.Context, log *logger.Logger, message *model.Message) error
	Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error
	Publish(ctx context.Context, log *logger.Logger, params *model.PublishMessagesParams) error
	Private(ctx context.Context, log *logger.Logger, params *model.PrivateMessagesParams) error
	Read(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (*model.Message, error)
	DeleteAllUserMessages(ctx context.Context, log *logger.Logger, params *model.DeleteAllUserMessagesParams) error
	ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (*model.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
	Truncate(ctx context.Context, log *logger.Logger) error
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	repo Repository
}

/**
 * Returns empty interface. It is either an error,
 * or new msg with unique id, saved in repo.
 * 
 * Apply replicates log state from the bottom up.
 * Leader makes Apply on start.
 */
func (f *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendRequest:
		return f.applyAppend(buf[1:])
	case UpdateRequest:
		return f.applyUpdate(buf[1:])
	case DeleteAllUserMessagesRequest:
		return f.applyDeleteAllUserMessages(buf[1:])
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case PublishRequest:
		return f.applyPublish(buf[1:])
	case PrivateRequest:
		return f.applyPrivate(buf[1:])
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	// Put does not put message with same id twice
	err := f.repo.Create(context.Background(), logger.Default(), model.MessageFromProto(cmd.Message))
	if err != nil {
		return err
	}
	return &AppendCommandResult{}
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	res, err := f.repo.Update(context.Background(), logger.Default(), &model.UpdateMessageParams{
		ID: cmd.Id,
		UserID: cmd.UserId,
		FileID: cmd.FileId,
		ThreadID: cmd.ThreadId,
		Text: cmd.Text,
		UpdateUTCNano: cmd.UpdateUtcNano,
		Private: cmd.Private,
	})
	if err != nil {
		return err
	}

	return &UpdateCommandResult{
		UpdateUtcNano: cmd.UpdateUtcNano,
		Private: res.Private,
	}
}

func (f *fsm) applyDeleteAllUserMessages(raw []byte) interface{} {
	var cmd DeleteAllUserMessagesCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.DeleteAllUserMessages(context.Background(), logger.Default(), &model.DeleteAllUserMessagesParams{
		UserID: cmd.UserId,
	})
	if err != nil {
		return err
	}

	return &DeleteCommandResult{}
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.Delete(context.Background(), logger.Default(), &model.DeleteMessageParams{
		ID: cmd.Id,
		UserID: cmd.UserId,
	})
	if err != nil {
		return err
	}

	return &DeleteCommandResult{Ok: true, Explain: "deleted"}
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.Publish(context.Background(), logger.Default(), &model.PublishMessagesParams{
		IDs: cmd.Ids,
		UserID: cmd.UserId,
		UpdateUTCNano: cmd.UpdateUtcNano,
	})
	if err != nil {
		return err
	}

	return &PublishCommandResult{}
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.Private(context.Background(), logger.Default(), &model.PrivateMessagesParams{
		IDs: cmd.Ids,
		UserID: cmd.UserId,
		UpdateUTCNano: cmd.UpdateUtcNano,
	})
	if err != nil {
		return err
	}

	return &PrivateCommandResult{}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

func (f *fsm) Restore(_ io.ReadCloser) error {
	return nil
}