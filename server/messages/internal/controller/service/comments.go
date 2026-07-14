package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/db/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
)

type CommentsConfig struct {
	RpcAddr string
}

type CommentsController struct {
	conf   CommentsConfig
	client comments.CommentsClient
	conn   *grpc.ClientConn
	publisher    ddd.EventPublisher[ddd.Event]
}

func NewCommentsController(conf CommentsConfig, publisher ddd.EventPublisher[ddd.Event]) *CommentsController {
	c := &CommentsController{
		conf: conf,
		publisher: publisher,
	}

	c.setupConnection()

	return c
}

func (s *CommentsController) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *CommentsController) setupConnection() (err error) {
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

	client := comments.NewCommentsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *CommentsController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		logger.Debugw("comments conn failed", "state", state.String())
		return true
	}
	return false
}

func (s *CommentsController) SendComment(ctx context.Context, id, userID, messageID int64, text string, metadata []byte) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("save comment", "id", id, "user_id", userID, "message_id", messageID, "text", text, "metadata", metadata)

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

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
		ReqType: int32(machine.AppendCommentRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.CreateComment(id, userID, messageID, text, createdAt, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *CommentsController) UpdateComment(ctx context.Context, id, userID int64, text *string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update comment", "id", id, "user_id", userID, "text", text)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

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
		ReqType: int32(machine.UpdateCommentRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.UpdateComment(id, userID, text, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *CommentsController) DeleteComment(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete comment", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&comments.DeleteCommentCommand{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteCommentRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.DeleteComment(id, userID)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *CommentsController) DeleteMessageComments(ctx context.Context, messageID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete message comments", "message_id", messageID)

	cmd, err := proto.Marshal(&comments.DeleteMessageCommentsCommand{
		MessageId: messageID,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteMessageCommentsRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.DeleteMessageComments(messageID)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *CommentsController) ReadComment(ctx context.Context, id, userID int64) (comment *model.Comment, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read comment", "id", id, "user_id", userID)

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

func (s *CommentsController) ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32, asc bool) (list *model.CommentsList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list comments", "message_id", messageID, "user_id", userID, "name", name, "limit", limit, "offset", offset, "ascending", asc)

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
