package balancer

import (
	"os"
	"fmt"
	"sync"
	"strings"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"github.com/bd878/gallery/server/internal/logger"
)

var _ base.PickerBuilder = (*SubchannelPicker)(nil)

type conn struct {
	SubConn  balancer.SubConn
	Address  string
}

type SubchannelPicker struct {
	mu               sync.RWMutex
	leader           *conn
	followers        []*conn
	current          uint64
	leaderMethods    []string
	followerMethods  []string
}

func (p *SubchannelPicker) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	p.mu.Lock()
	defer p.mu.Unlock()

	logger.Debugln("build picker")

	p.leader = nil

	var followers []*conn
	for sc, scInfo := range buildInfo.ReadySCs {
		logger.Debugw("scinfo", "address", scInfo.Address.Addr, "is_leader", scInfo.Address.Attributes.Value("is_leader"))

		isLeader := scInfo.
			Address.
			Attributes.
			Value("is_leader").(bool)

		sc.RegisterHealthListener(func (subConn balancer.SubConnState) {
			fmt.Fprintln(os.Stdout, "address", scInfo.Address.Addr,
				"ConnectivityState", subConn.ConnectivityState.String(), "ConnectionError", subConn.ConnectionError)
		})

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

var _ balancer.Picker = (*SubchannelPicker)(nil)

func (p *SubchannelPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	logger.Debugw("picker pick", "FullMethodName", info.FullMethodName)

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

	result.Done = func(info balancer.DoneInfo) {
		fmt.Fprintf(os.Stdout, "done info [%T]: %+[1]v\n", info)
	}

	fmt.Fprintln(os.Stdout, "pick conn", "method_name", info.FullMethodName, "addr", endpoint.Address)
	return result, nil
}

func (p *SubchannelPicker) nextFollower() *conn {
	cur := atomic.AddUint64(&p.current, uint64(1))
	len := uint64(len(p.followers))
	idx := int(cur % len)
	return p.followers[idx]
}

func RegisterPicker(name string, leaderMethods []string, followerMethods []string) {
	balancer.Register(
		base.NewBalancerBuilder(name, &SubchannelPicker{leaderMethods: leaderMethods, followerMethods: followerMethods}, base.Config{HealthCheck: true}),
	)
}
