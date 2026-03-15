package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
)

type CommentsController interface {
	SaveComment(ctx context.Context, id, userID, messageID int64, text string, metadata []byte) (err error)
	DeleteComment(ctx context.Context, id, userID int64) (err error)
	DeleteMessageComments(ctx context.Context, messageID int64) (err error)
	UpdateComment(ctx context.Context, id, userID int64, text *string, metadata []byte) (err error)
	ReadComment(ctx context.Context, id, userID int64) (comment *api.Comment, err error)
	ListComments(ctx context.Context, messageID, userID *int64, limit, offset int32) (list *api.CommentsList, err error)
}

type CommentsHandler struct {
	api.UnimplementedCommentsServer
	controller CommentsController
}

func NewCommentsHandler(ctrl CommentsController) *CommentsHandler {
	return &CommentsHandler{controller: ctrl}
}

func (h *CommentsHandler) SendComment(ctx context.Context, req *api.SendCommentRequest) (resp *api.SendCommentResponse, err error) {
	err = h.controller.SaveComment(ctx, req.Id, req.UserId, req.MessageId, req.Text, req.Metadata)

	resp = &api.SendCommentResponse{
	}

	return
}

func (h *CommentsHandler) UpdateComment(ctx context.Context, req *api.UpdateCommentRequest) (resp *api.UpdateCommentResponse, err error) {
	err = h.controller.UpdateComment(ctx, req.Id, req.UserId, req.Text, req.Metadata)

	resp = &api.UpdateCommentResponse{}

	return
}

func (h *CommentsHandler) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (resp *api.DeleteCommentResponse, err error) {
	err = h.controller.DeleteComment(ctx, req.Id, req.UserId)

	resp = &api.DeleteCommentResponse{}

	return
}

func (h *CommentsHandler) DeleteMessageComments(ctx context.Context, req *api.DeleteMessageCommentsRequest) (resp *api.DeleteMessageCommentsResponse, err error) {
	err = h.controller.DeleteMessageComments(ctx, req.MessageId)

	resp = &api.DeleteMessageCommentsResponse{}

	return
}

func (h *CommentsHandler) ReadComment(ctx context.Context, req *api.ReadCommentRequest) (resp *api.ReadCommentResponse, err error) {
	comment, err := h.controller.ReadComment(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadCommentResponse{
		Comment: comment,
	}

	return
}

func (h *CommentsHandler) ListComments(ctx context.Context, req *api.ListCommentsRequest) (resp *api.ListCommentsResponse, err error) {
	list, err := h.controller.ListComments(ctx, req.MessageId, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	resp = &api.ListCommentsResponse{
		Comments:    list.Comments,
		IsLastPage:  list.IsLastPage,
		IsFirstPage: list.IsFirstPage,
		Total:       list.Total,
		Count:       list.Count,
	}

	return
}
