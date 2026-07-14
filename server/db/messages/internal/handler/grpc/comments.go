package grpc

import (
	"context"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
)

type CommentsController interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	ReadComment(ctx context.Context, id, userID int64) (comment *comments.Comment, err error)
	ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32) (list *comments.CommentsList, err error)
}

type CommentsHandler struct {
	comments.UnimplementedCommentsServer
	controller CommentsController
}

func NewCommentsHandler(ctrl CommentsController) *CommentsHandler {
	return &CommentsHandler{controller: ctrl}
}

func (h *CommentsHandler) Apply(ctx context.Context, req *api.Command) (resp *api.CommandResponse, err error) {
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	err = h.controller.Apply(ctx, machine.RequestType(req.ReqType), req.Cmd, duration)

	resp = &api.CommandResponse{}

	return
}

func (h *CommentsHandler) ReadComment(ctx context.Context, req *comments.ReadCommentRequest) (resp *comments.ReadCommentResponse, err error) {
	comment, err := h.controller.ReadComment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &comments.ReadCommentResponse{
		Comment: comment,
	}

	return
}

func (h *CommentsHandler) ListComments(ctx context.Context, req *comments.ListCommentsRequest) (resp *comments.ListCommentsResponse, err error) {
	list, err := h.controller.ListComments(ctx, req.UserId, req.MessageId, req.Name, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	resp = &comments.ListCommentsResponse{
		Comments:    list.Comments,
		IsLastPage:  list.IsLastPage,
		IsFirstPage: list.IsFirstPage,
		Total:       list.Total,
		Count:       list.Count,
	}

	return
}
