package distributed

import (
	"time"
	"fmt"
	"context"
	"errors"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

func (m *DistributedMessages) apply(ctx context.Context, reqType RequestType, cmd []byte) (
	interface{}, error,
) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return nil, err
	}

	timeout := 10*time.Second
	/* fsm.Apply() */
	future := m.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res := future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res, nil
}

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) error {
	cmd, _ := proto.Marshal(&AppendCommand{
		Message: model.MessageToProto(message),
	})

	_, err := m.apply(ctx, AppendRequest, cmd)
	if err != nil {
		log.Errorln("raft failed to apply save message")
		return err
	}

	return nil
}

func (m *DistributedMessages) UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (
	*model.UpdateMessageResult, error,
) {
	cmd, _ := proto.Marshal(&UpdateCommand{
		Id: params.ID,
		UserId: params.UserID,
		FileId: params.FileID,
		ThreadId: params.ThreadID,
		Text:   params.Text,
		UpdateUtcNano: params.UpdateUTCNano,
		Private: params.Private,
	})

	res, err := m.apply(ctx, UpdateRequest, cmd)
	if err != nil {
		log.Errorln("raft failed to apply save message")
		return nil, err
	}

	switch val := res.(type) {
	case *UpdateCommandResult:
		return &model.UpdateMessageResult{
			Private: val.Private,
		}, nil
	case error:
		return nil, val
	default:
		log.Errorln("update request reseived unknown type")
		return nil, errors.New("unknown message update type")
	}
}

func (m *DistributedMessages) DeleteAllUserMessages(ctx context.Context, log *logger.Logger, params *model.DeleteAllUserMessagesParams) error {
	cmd, _ := proto.Marshal(&DeleteAllUserMessagesCommand{
		UserId: params.UserID,
	})

	_, err := m.apply(ctx, DeleteAllUserMessagesRequest, cmd)
	if err != nil {
		log.Errorln("raft failed to apply delete all messages")
		return err
	}

	return nil
}

func (m *DistributedMessages) DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
	cmd, _ := proto.Marshal(&DeleteCommand{
		Id: params.ID,
		UserId: params.UserID,
	})

	_, err := m.apply(ctx, DeleteRequest, cmd)
	if err != nil {
		log.Errorln("raft failed to apply delete message")
		return err
	}

	return nil
}

func (m *DistributedMessages) DeleteMessages(ctx context.Context, log *logger.Logger, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error) {
	statuses := make([]*model.DeleteMessageStatus, 0, len(params.IDs))
	for _, id := range params.IDs {
		cmd, _ := proto.Marshal(&DeleteCommand{
			Id: id,
			UserId: params.UserID,
		})

		res, err := m.apply(ctx, DeleteRequest, cmd)
		if err != nil {
			log.Errorw("raft failed to apply delete message", "error", err)
			statuses = append(statuses, &model.DeleteMessageStatus{
				ID: id,
				OK: false,
				Explain: "error",
			})
			continue
		}

		status, ok := res.(*DeleteCommandResult)
		if !ok {
			return nil, fmt.Errorf("cannot cast %T to *DeleteCommandResult\n", status)
		}

		statuses = append(statuses, &model.DeleteMessageStatus{
			ID: id,
			OK: status.Ok,
			Explain: status.Explain,
		})
	}

	return &model.DeleteMessagesResult{IDs: statuses}, nil
}

func (m *DistributedMessages) PublishMessages(ctx context.Context, log *logger.Logger, params *model.PublishMessagesParams) (*model.PublishMessagesResult, error) {
	cmd, _ := proto.Marshal(&PublishCommand{
		Ids: params.IDs,
		UserId: params.UserID,
		UpdateUtcNano: params.UpdateUTCNano,
	})

	_, err := m.apply(ctx, PublishRequest, cmd)
	if err != nil {
		return nil, err
	}

	return &model.PublishMessagesResult{
		UpdateUTCNano: params.UpdateUTCNano,
	}, nil
}

func (m *DistributedMessages) PrivateMessages(ctx context.Context, log *logger.Logger, params *model.PrivateMessagesParams) (*model.PrivateMessagesResult, error) {
	cmd, _ := proto.Marshal(&PrivateCommand{
		Ids: params.IDs,
		UserId: params.UserID,
		UpdateUtcNano: params.UpdateUTCNano,
	})

	_, err := m.apply(ctx, PrivateRequest, cmd)
	if err != nil {
		return nil, err
	}

	return &model.PrivateMessagesResult{
		UpdateUTCNano: params.UpdateUTCNano,
	}, nil
}

func (m *DistributedMessages) ReadMessage(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (
	*model.Message, error,
) {
	return m.repo.Read(ctx, log, params)
}

func (m *DistributedMessages) ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (
	*model.ReadThreadMessagesResult, error,
) {
	log.Infow("distributed methods. read thread messages", "user_id", params.UserID, "thread_id", params.ThreadID)
	return m.repo.ReadThreadMessages(
		ctx,
		log,
		&model.ReadThreadMessagesParams{
			UserID:    params.UserID,
			ThreadID:  params.ThreadID,
			Limit:     params.Limit,
			Offset:    params.Offset,
			Ascending: params.Ascending,
			Private:   params.Private,
		},
	)
}

func (m *DistributedMessages) ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (
	*model.ReadMessagesResult, error,
) {
	return m.repo.ReadAllMessages(
		ctx,
		log,
		&model.ReadMessagesParams{
			UserID:    params.UserID,
			Limit:     params.Limit,
			Offset:    params.Offset,
			Ascending: params.Ascending,
			Private:   params.Private,
		},
	)
}
