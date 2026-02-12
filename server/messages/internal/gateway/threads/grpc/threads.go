package grpc

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/threads/pkg/loadbalance"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Gateway struct {
	addr    string
	client  api.ThreadsClient
	conn    *grpc.ClientConn
}

func New(addr string) *Gateway {
	gateway := &Gateway{addr: addr}

	gateway.setupConnection()

	return gateway
}

func (g *Gateway) setupConnection() error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			g.addr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = api.NewThreadsClient(conn)

	return nil
}

func (g *Gateway) isConnFailed() bool {
	state := g.conn.GetState()
	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func (g *Gateway) CreateThread(ctx context.Context, id, userID, parentID int64, name string, private bool) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create thread", "id", id, "user_id", userID, "parent_id", parentID, "name", name, "private", private)

	_, err = g.client.Create(ctx, &api.CreateRequest{
		Id:        id,
		UserId:    userID,
		ParentId:  parentID,
		Name:      name,
		Private:   private,
	})

	return
}

func (g *Gateway) DeleteThread(ctx context.Context, id, userID int64) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete thread", "id", id, "user_id", userID)

	_, err = g.client.Delete(ctx, &api.DeleteRequest{
		Id:     id,
		UserId: userID,
	})

	return
}

func (g *Gateway) UpdateThread(ctx context.Context, id, userID int64) (err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update thread", "id", id, "user_id", userID)

	_, err = g.client.Update(ctx, &api.UpdateRequest{
		Id:        id,
		UserId:    userID,
	})

	return
}

func (g *Gateway) ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32) (list []*threads.Thread, isLastPage bool, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list threads", "user_id", userID, "parent_id", parentID, "limit", limit, "offset", offset)

	resp, err := g.client.List(ctx, &api.ListRequest{
		UserId:   userID,
		ParentId: parentID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, false, err
	}

	isLastPage = resp.IsLastPage
	list = threads.MapThreadsFromProto(threads.ThreadFromProto, resp.List)

	return
}

func (g *Gateway) CountThreads(ctx context.Context, id, userID int64) (total int32, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("count threads", "id", id, "user_id", userID)

	resp, err := g.client.Count(ctx, &api.CountRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return 0, err
	}

	total = resp.Total

	return
}

func (g *Gateway) ResolvePath(ctx context.Context, userID, id int64) (path []int64, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("resolve path", "user_id", userID, "id", id)

	resp, err := g.client.Resolve(ctx, &api.ResolveRequest{
		UserId:  userID,
		Id:      id,
	})
	if err != nil {
		return nil, err
	}

	path = resp.Path

	return
}

func (g *Gateway) ReadThread(ctx context.Context, userID, id int64) (thread *threads.Thread, err error) {
	if g.isConnFailed() {
		if err = g.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read thread", "user_id", userID, "id", id)

	resp, err := g.client.Read(ctx, &api.ReadRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	thread = threads.ThreadFromProto(resp)

	return
}