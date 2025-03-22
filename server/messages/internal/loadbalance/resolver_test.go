package loadbalance_test

import (
	"net"
	"net/url"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/attributes"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/internal/loadbalance"
)

func TestResolver(t *testing.T) {
	serv := NewGRPCServer()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go serv.Serve(l)

	conn := &clientConn{}
	r := &loadbalance.Resolver{}
	_, err = r.Build(
		resolver.Target{
			URL: url.URL{
				Scheme: "messages",
				Host: "",
				Path: l.Addr().String(),
			},
		},
		conn,
		resolver.BuildOptions{},
	)
	require.NoError(t, err)

	wantState := resolver.State{
		Endpoints: []resolver.Endpoint{{
			Addresses: []resolver.Address{{
				Addr: "localhost:9001",
				Attributes: attributes.New("is_leader", true),
			}, {
				Addr: "localhost:9002",
				Attributes: attributes.New("is_leader", false),
			}},
		}},
	}
	require.Equal(t, wantState, conn.state)
	conn.state.Addresses = nil
	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

type clientConn struct {
	state resolver.State
}

func (cc *clientConn) UpdateState(state resolver.State) error {
	cc.state = state
	return nil
}

func (cc *clientConn) ReportError(err error) {}

func (cc *clientConn) NewAddress([]resolver.Address) {}

func (cc *clientConn) NewServiceConfig(config string) {}

func (cc *clientConn) ParseServiceConfig(string) *serviceconfig.ParseResult {
	return nil
}

type grpcServer struct {
	api.UnimplementedMessagesServiceServer
}

func NewGRPCServer() *grpc.Server {
	gsrv := grpc.NewServer()

	srv := &grpcServer{}

	api.RegisterMessagesServiceServer(gsrv, srv)
	return gsrv
}

func (s *grpcServer) GetServers(_ context.Context, _ *api.GetServersRequest) (
	*api.GetServersResponse, error,
) {
	servers := []*api.Server{{
		Id: "leader",
		RaftAddr: "localhost:9001",
		IsLeader: true,
	}, {
		Id: "follower",
		RaftAddr: "localhost:9002",
	}}

	return &api.GetServersResponse{Servers: servers}, nil
}

func (s *grpcServer) SaveMessage(_ context.Context, _ *api.SendMessageRequest) (
	*api.SaveMessageResponse, error,
) {
	return nil, nil
}

func (s *grpcServer) ReadAllMessages(_ context.Context, _ *api.ReadMessageRequest) (
	*api.ReadMessageResponse, error,
) {
	return nil, nil
}