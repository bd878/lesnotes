package loadbalance

import (
  "context"
  "sync"
  "fmt"
  "log"

  "google.golang.org/grpc"
  "google.golang.org/grpc/attributes"
  "google.golang.org/grpc/resolver"
  "google.golang.org/grpc/serviceconfig"

  "github.com/bd878/gallery/server/gen"
)

type Resolver struct {
  mu sync.Mutex
  clientConn resolver.ClientConn
  resolverConn *grpc.ClientConn
  serviceConfig *serviceconfig.ParseResult
}

var _ resolver.Builder = (*Resolver)(nil)

func (r *Resolver) Build(
  target resolver.Target,
  cc resolver.ClientConn,
  opts resolver.BuildOptions,
) (resolver.Resolver, error) {
  var err error
  r.serviceConfig = r.clientConn.ParseServiceConfig(
    fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, Name),
  )

  var dialOpts []grpc.DialOption
  dialOpts = append(dialOpts, grpc.WithTransportCredentials(opts.DialCreds))
  r.resolverConn, err = grpc.Dial(target.Endpoint(), dialOpts...)
  if err != nil {
    return nil, err
  }
  r.ResolveNow(resolver.ResolveNowOptions{})
  return r, nil
}

const Name = "messages"

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

  client := gen.NewMessagesServiceClient(r.resolverConn)

  ctx := context.Background()
  res, err := client.GetServers(ctx, &gen.GetMessagesServersRequest{})
  if err != nil {
    return
  }

  var addrs []resolver.Address
  for _, server := range res.Servers {
    addrs = append(addrs, resolver.Address{
      Addr: server.RpcAddr,
      Attributes: attributes.New(
        "is_leader",
        server.IsLeader,
      ),
    })
  }
  r.clientConn.UpdateState(resolver.State{
    Endpoints: []resolver.Endpoint{{
      Addresses: addrs,
    }},
    ServiceConfig: r.serviceConfig,
  })
}

func (r *Resolver) Close() {
  if err := r.resolverConn.Close(); err != nil {
    log.Println(err)
  }
}