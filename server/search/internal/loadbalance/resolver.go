package loadbalance

import (
	"context"
	"sync"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/resolver"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
)

type Resolver struct {
	mu             sync.Mutex
	clientConn     resolver.ClientConn
	resolverConn   *grpc.ClientConn
	serviceConfig  *serviceconfig.ParseResult
}

var _ resolver.Builder = (*Resolver)(nil)

func (r *Resolver) Build(
	t resolver.Target,
	cc resolver.ClientConn,
	_ resolver.BuildOptions,
) (resolver.Resolver, error) {
	var err error

	logger.Debugw("build resolver", "endpoing", t.Endpoint())

	r.clientConn = cc
	r.resolverConn, err = grpc.NewClient(
		t.Endpoint(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	r.serviceConfig = r.clientConn.ParseServiceConfig(
		fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, Name),
	)
	if err != nil {
		return nil, err
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

const Name = "search"

func (r *Resolver) Scheme() string {
	return Name
}

func init() {
	resolver.Register(&Resolver{})
}

var _ resolver.Resolver = (*Resolver)(nil)

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client := api.NewSearchClient(r.resolverConn)

	logger.Debugw("resolver resolve now servers", "client", r.resolverConn.Target())

	ctx := context.Background()
	res, err := client.GetServers(ctx, &api.GetServersRequest{})
	if err != nil {
		logger.Errorw("failed to get servers", "error", err)
		return
	}

	logger.Debugw("resolver received new servers", "servers", res.Servers)

	var addrs []resolver.Address
	for _, server := range res.Servers {
		addrs = append(addrs, resolver.Address{
			Addr: server.RaftAddr,
			Attributes: attributes.New(
				"is_leader",
				server.IsLeader,
			),
		})
	}
	r.clientConn.UpdateState(resolver.State{
		Addresses: addrs,
		ServiceConfig: r.serviceConfig,
	})
}

func (r *Resolver) Close() {
	if err := r.resolverConn.Close(); err != nil {
		logger.Error(err)
	}
}