package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
)

type CommentsConfig struct {
	RpcAddr string
}

type CommentsController struct {
	conf       CommentsConfig
	client     api.CommentsClient
	conn       *grpc.ClientConn
}

func NewCommentsController(conf CommentsConfig) *CommentsController {
	c := &CommentsController{
		conf: conf,
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

	client := api.NewCommentsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *CommentsController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugw("connection failed", "state", state.String())
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

	_, err = s.client.SendComment(ctx, &api.SendCommentRequest{
		Id:         id,
		UserId:     userID,
		MessageId:  messageID,
		Text:       text,
		Metadata:   metadata,
	})

	return
}

func (s *CommentsController) UpdateComment(ctx context.Context, id, userID int64, text *string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update comment", "id", id, "user_id", userID, "text", text)

	_, err = s.client.UpdateComment(ctx, &api.UpdateCommentRequest{
		Id:       id,
		UserId:   userID,
		Text:     text,
	})

	return
}

func (s *CommentsController) DeleteComment(ctx context.Context, id, userID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete comment", "id", id, "user_id", userID)

	_, err = s.client.DeleteComment(ctx, &api.DeleteCommentRequest{
		Id:     id,
		UserId: userID,
	})

	return
}

func (s *CommentsController) DeleteMessageComments(ctx context.Context, messageID int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete message comments", "message_id", messageID)

	_, err = s.client.DeleteMessageComments(ctx, &api.DeleteMessageCommentsRequest{
		MessageId: messageID,
	})

	return
}

func (s *CommentsController) ReadComment(ctx context.Context, id, userID int64) (comment *model.Comment, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read comment", "id", id, "user_id", userID)

	res, err := s.client.ReadComment(ctx, &api.ReadCommentRequest{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return
	}

	comment = model.CommentFromProto(res.Comment)

	return
}

func (s *CommentsController) ListComments(ctx context.Context, userID, messageID *int64, limit, offset int32, asc bool) (list *model.CommentsList, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list comments", "message_id", messageID, "user_id", userID, "limit", limit, "offset", offset, "ascending", asc)

	res, err := s.client.ListComments(ctx, &api.ListCommentsRequest{
		MessageId:     messageID,
		UserId:        userID,
		Limit:         limit,
		Offset:        offset,
		Asc:           asc,
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
