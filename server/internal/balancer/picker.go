package balancer

import (
	"os"
	"fmt"
	"sync"
	"strings"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

var _ base.PickerBuilder = (*Picker)(nil)

type conn struct {
	SubConn  balancer.SubConn
	Address  string
}

type Picker struct {
	mu               sync.RWMutex
	leader           *conn
	followers        []*conn
	current          uint64
	leaderMethods    []string
	followerMethods  []string
}

func (p *Picker) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	p.mu.Lock()
	defer p.mu.Unlock()

	var followers []*conn
	for sc, scInfo := range buildInfo.ReadySCs {
		isLeader := scInfo.
			Address.
			Attributes.
			Value("is_leader").(bool)

		if isLeader {
			p.leader = &conn{SubConn: sc, Address: scInfo.Address.Addr}
			continue
		}
		followers = append(followers, &conn{SubConn: sc, Address: scInfo.Address.Addr})
	}
	p.followers = followers
	fmt.Fprintln(os.Stdout, "len(followers)", len(p.followers))
	return p
}

var _ balancer.Picker = (*Picker)(nil)

func (p *Picker) Pick(info balancer.PickInfo) (
	balancer.PickResult, error,
) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result balancer.PickResult

	endpoint := p.leader

	if len(p.followers) > 0 {
		for _, method := range p.followerMethods {
			if strings.Contains(info.FullMethodName, method) {
				endpoint = p.nextFollower()
			}
		}
	}

	if endpoint == nil {
		fmt.Fprintln(os.Stdout, "no sub conn available", "method_name", info.FullMethodName)
		return result, balancer.ErrNoSubConnAvailable
	}
	result.SubConn = endpoint.SubConn
	fmt.Fprintln(os.Stdout, "pick conn", "method_name", info.FullMethodName, "addr", endpoint.Address)
	return result, nil
}

func (p *Picker) nextFollower() *conn {
	cur := atomic.AddUint64(&p.current, uint64(1))
	len := uint64(len(p.followers))
	idx := int(cur % len)
	return p.followers[idx]
}

func RegisterPicker(name string, leaderMethods []string, followerMethods []string) {
	balancer.Register(
		base.NewBalancerBuilder(name, &Picker{leaderMethods: leaderMethods, followerMethods: followerMethods}, base.Config{}),
	)
}
