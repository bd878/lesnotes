package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/grpc/status"

	"github.com/bd878/gallery/server/api"
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

	slog.Debug("build resolver", slog.String("endpoint", r.target))

	r.clientConn = cc
	r.resolverConn, err = grpc.NewClient(
		r.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	r.serviceConfig = r.clientConn.ParseServiceConfig(
		fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, r.name),
	)
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

		slog.Debug("resolver resolve now servers", slog.Int("retry", i), slog.String("target", r.resolverConn.Target()))

		res, err := client.GetServers(context.TODO(), &api.GetServersRequest{})
		if err != nil {
			if status, ok := status.FromError(err); ok {
				slog.Debug("failed to get servers", slog.String("code", status.Code().String()), slog.String("status", status.Message()))

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

					if err := r.resolverConn.Close(); err != nil {
						slog.Error("failed to close old resolver conn", slog.String("error", err.Error()))
					}

					r.resolverConn, err = grpc.NewClient(
						r.target,
						grpc.WithTransportCredentials(insecure.NewCredentials()),
					)
					if err != nil {
						slog.Error("failed to recreate resolver conn", slog.String("error", err.Error()))
						return
					}
					r.serviceConfig = r.clientConn.ParseServiceConfig(
						fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, r.name),
					)

					slog.Debug("recover", slog.String("target", r.target))

					continue
				}
			}

			slog.Error("failed to get servers", slog.String("error", err.Error()))
			r.clientConn.ReportError(err)

			if err := r.resolverConn.Close(); err != nil {
				slog.Error("failed to close old resolver conn", slog.String("error", err.Error()))
			}

			return
		}

		slog.Debug("resolver received new servers", slog.String("servers", fmt.Sprintf("%v", res.Servers)))

		var addrs []resolver.Address
		var hasLeader bool
		for _, server := range res.Servers {
			slog.Debug("server", slog.String("raft_addr", server.RaftAddr), slog.String("is_leader", fmt.Sprintf("%v", server.IsLeader)))

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

		slog.Debug("no leader")

		<-ticker.C
	}
}

func (r *Resolver) Close() {
	if r.resolverConn != nil {
		slog.Debug("close resolver conn")
		if err := r.resolverConn.Close(); err != nil {
			slog.Error("close resolver conn", slog.String("error", err.Error()))
		}
	}
}
