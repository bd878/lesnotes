package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/search/pkg/loadbalance"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

type Config struct {
	RpcAddr  string
}

type Controller struct {
	conf         Config
	client       api.SearchClient
	conn         *grpc.ClientConn
}

func New(conf Config) *Controller {
	controller := &Controller{conf: conf}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Controller) setupConnection() (err error) {
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

	client := api.NewSearchClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
		return true
	}
	return false
}

func (s *Controller) SearchMessages(ctx context.Context, userID int64, substr string) (list []*searchmodel.Message, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("search messages", "user_id", userID, "substr", substr)

	res, err := s.client.SearchMessages(ctx, &api.SearchMessagesRequest{
		Substr: substr,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	list = searchmodel.MapMessagesFromProto(searchmodel.MessageFromProto, res.List)

	return 
}