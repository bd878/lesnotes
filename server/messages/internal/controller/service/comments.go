package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type CommentsController struct {
	client    comments.CommentsClient
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCommentsController(container di.Container, publisher ddd.EventPublisher[ddd.Event]) CommentsController {
	client := container.Get("commentsClient").(comments.CommentsClient)

	return CommentsController{
		client: client,
		publisher: publisher,
	}
}

func (s CommentsController) SendComment(ctx context.Context, id, userID, messageID int64, text string, metadata []byte) (err error) {
	slog.Debug("save comment", slog.Int64("id", id), slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("text", text), slog.String("metadata", fmt.Sprintf("%v", metadata)))

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.CreateComment(id, userID, messageID, text, createdAt, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&comments.AppendCommentCommand{
		Id:        id,
		UserId:    userID,
		MessageId: messageID,
		Text:      text,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendCommentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s CommentsController) UpdateComment(ctx context.Context, id, userID int64, text *string) (err error) {

	logValues := []any{slog.Int64("id", id), slog.Int64("user_id", userID)}
	if text != nil {
		logValues = append(logValues, slog.String("text", *text))
	}
	slog.Debug("update comment", logValues...)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateComment(id, userID, text, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&comments.UpdateCommentCommand{
		Id:        id,
		UserId:    userID,
		Text:      text,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.UpdateCommentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s CommentsController) DeleteComment(ctx context.Context, id, userID int64) (err error) {
	slog.Debug("delete comment", slog.Int64("id", id), slog.Int64("user_id", userID))

	event, err := domain.DeleteComment(id, userID)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&comments.DeleteCommentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.DeleteCommentRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return 
}

func (s CommentsController) DeleteMessageComments(ctx context.Context, messageID int64) (err error) {
	slog.Debug("delete message comments", slog.Int64("message_id", messageID))

	event, err := domain.DeleteMessageComments(messageID)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&comments.DeleteMessageCommentsCommand{
		MessageId: messageID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.DeleteMessageCommentsRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s CommentsController) ReadComment(ctx context.Context, id, userID int64) (comment *model.Comment, err error) {
	slog.Debug("read comment", slog.Int64("id", id), slog.Int64("user_id", userID))

	res, err := s.client.ReadComment(ctx, &comments.ReadCommentRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return
	}

	comment = model.CommentFromProto(res.Comment)

	return
}

func (s CommentsController) ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32, asc bool) (list *model.CommentsList, err error) {

	logValues := []any{slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("ascending", asc)}
	if messageID != nil {
		logValues = append(logValues, slog.Int64("message_id", *messageID))
	}
	if userID != nil {
		logValues = append(logValues, slog.Int64("user_id", *userID))
	}
	if name != nil {
		logValues = append(logValues, slog.String("name", *name))
	}
	slog.Debug("list comments", logValues...)

	res, err := s.client.ListComments(ctx, &comments.ListCommentsRequest{
		MessageId: messageID,
		UserId:    userID,
		Name:      name,
		Limit:     limit,
		Offset:    offset,
		Asc:       asc,
	})
	if err != nil {
		return
	}

	list = &model.CommentsList{
		Comments:    model.MapCommentsFromProto(model.CommentFromProto, res.Comments),
		IsLastPage:  res.IsLastPage,
		IsFirstPage: res.IsFirstPage,
		Total:       res.Total,
		Count:       res.Count,
		Offset:      offset,
	}

	return
}
