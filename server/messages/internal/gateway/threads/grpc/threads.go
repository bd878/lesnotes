package grpc

import (
	"context"

	"log/slog"

	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/threads/pkg/model"
)

type Gateway struct {
	client threads.ThreadsClient
}

func New(container di.Container) *Gateway {
	client := container.Get("threadsClient").(threads.ThreadsClient)

	return &Gateway{client}
}

func (g *Gateway) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32) (list []*model.Thread, isLastPage bool, err error) {
	slog.Debug("list threads", slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)))

	resp, err := g.client.List(ctx, &threads.ListRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, false, err
	}

	isLastPage = resp.IsLastPage
	list = model.MapThreadsFromProto(model.ThreadFromProto, resp.List)

	return
}

func (g *Gateway) ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, privateMessage *bool) (list []*model.Thread, isLastPage bool, err error) {

	logValues := []any{slog.Int64("user_id", userID), slog.Int64("parent_id", parentID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset))}
	if privateMessage != nil {
		logValues = append(logValues, slog.Bool("private_message", *privateMessage))
	}
	slog.Debug("list messages", logValues...)

	resp, err := g.client.ListMessages(ctx, &threads.ListMessagesRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
		Private:  privateMessage,
	})
	if err != nil {
		return nil, false, err
	}

	isLastPage = resp.IsLastPage
	list = model.MapThreadsFromProto(model.ThreadFromProto, resp.List)

	return
}

func (g *Gateway) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	slog.Debug("count threads", slog.Int64("id", id), slog.Int64("user_id", userID))

	resp, err := g.client.Count(ctx, &threads.CountRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return 0, err
	}

	total = resp.Total

	return
}

func (g *Gateway) CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error) {

	logValues := []any{slog.Int64("id", id), slog.Int64("user_id", userID)}
	if privateMessage != nil {
		logValues = append(logValues, slog.Bool("private_message", *privateMessage))
	}

	slog.Debug("count messages", logValues...)

	resp, err := g.client.CountMessages(ctx, &threads.CountMessagesRequest{
		UserId:         userID,
		Id:             id,
		PrivateMessage: privateMessage,
	})
	if err != nil {
		return 0, err
	}

	total = resp.Total

	return
}

func (g *Gateway) ResolvePath(ctx context.Context, userID, id int64) (path []*threads.PathStep, err error) {
	slog.Debug("resolve path", slog.Int64("user_id", userID), slog.Int64("id", id))

	resp, err := g.client.Resolve(ctx, &threads.ResolveRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return nil, err
	}

	path = resp.Path

	return
}

func (g *Gateway) ReadThread(ctx context.Context, userID, id int64, name string) (thread *model.Thread, err error) {
	slog.Debug("read thread", slog.Int64("user_id", userID), slog.Int64("id", id), slog.String("name", name))

	resp, err := g.client.Read(ctx, &threads.ReadRequest{
		Id:     id,
		UserId: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	thread = model.ThreadFromProto(resp)

	return
}

func (g *Gateway) ReadParent(ctx context.Context, userID, id int64) (thread *model.Thread, err error) {
	slog.Debug("read parent", slog.Int64("user_id", userID), slog.Int64("id", id))

	resp, err := g.client.ReadParent(ctx, &threads.ReadParentRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return nil, err
	}

	if resp.IsRoot {
		return &model.Thread{
			ID:      0,
			Name:    "",
			Private: true,
		}, nil
	}

	thread = model.ThreadFromProto(resp.Parent)

	return
}
