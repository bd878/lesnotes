package balancer

import (
	"time"
	"context"
	"sync"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/resolver"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Resolver struct {
	name           string
	mu             sync.Mutex
	addrs          []resolver.Address
	target         string
	clientConn     resolver.ClientConn
	resolverConn   *grpc.ClientConn
	serviceConfig  *serviceconfig.ParseResult
}

var _ resolver.Builder = (*Resolver)(nil)

func (r *Resolver) Build(t resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	var err error

	// primary target may become unavailable
	if r.target == "" {
		r.target = t.Endpoint()
	}

	logger.Debugw("build resolver", "endpoint", r.target)

	r.clientConn = cc
	r.resolverConn, err = grpc.NewClient(
		r.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	r.serviceConfig = r.clientConn.ParseServiceConfig(
		fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, r.name),
	)
	if err != nil {
		return nil, err
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.name
}

func RegisterResolver(name string) {
	resolver.Register(&Resolver{name: name})
}

var _ resolver.Resolver = (*Resolver)(nil)

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticker := time.NewTicker(1*time.Second)
	defer ticker.Stop()
	for i := range 10 {
		client := api.NewDistributedClient(r.resolverConn)

		logger.Debugw("resolver resolve now servers", "retry", i, "target", r.resolverConn.Target())

		res, err := client.GetServers(context.TODO(), &api.GetServersRequest{})
		if err != nil {
			if status, ok := status.FromError(err); ok {
				logger.Debugw("failed to get servers", "code", status.Code().String(), "status", status.Message())

				if len(r.addrs) > 0 && (
					status.Code() == codes.Unavailable ||
					status.Code() == codes.DeadlineExceeded ||
					status.Code() == codes.Canceled ||
					status.Code() == codes.NotFound ) {

					var addrs []resolver.Address
					for _, addr := range r.addrs {
						if addr.Addr != r.resolverConn.Target() {
							addrs = append(addrs, addr)
							r.target = addr.Addr
						}
					}

					r.addrs = addrs

					r.resolverConn, err = grpc.NewClient(
						r.target,
						grpc.WithTransportCredentials(insecure.NewCredentials()),
					)
					r.serviceConfig = r.clientConn.ParseServiceConfig(
						fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, r.name),
					)

					logger.Debugw("recover", "target", r.target)

					continue
				}
			}

			logger.Errorw("failed to get servers", "error", err)
			r.clientConn.ReportError(err)

			return
		}

		logger.Debugw("resolver received new servers", "servers", res.Servers)

		var addrs []resolver.Address
		var hasLeader bool
		for _, server := range res.Servers {
			logger.Debugw("server", "raft_addr", server.RaftAddr, "is_leader", server.IsLeader)

			if server.IsLeader {
				hasLeader = true
			}

			addrs = append(addrs, resolver.Address{
				Addr: server.RaftAddr,
				Attributes: attributes.New(
					"is_leader",
					server.IsLeader,
				),
			})
		}

		if hasLeader {
			r.addrs = addrs
			r.clientConn.UpdateState(resolver.State{
				Addresses: addrs,
				ServiceConfig: r.serviceConfig,
			})

			return
		}

		logger.Debugln("no leader")

		<-ticker.C
	}
}

func (r *Resolver) Close() {
	if err := r.resolverConn.Close(); err != nil {
		logger.Error(err)
	}
}