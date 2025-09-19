package loadbalance

import (
	"sync"
	"strings"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

var _ base.PickerBuilder = (*Picker)(nil)

type Picker struct {
	mu sync.RWMutex
	leader balancer.SubConn
	followers []balancer.SubConn
	current uint64
}

func (p *Picker) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	p.mu.Lock()
	defer p.mu.Unlock()

	var followers []balancer.SubConn
	for sc, scInfo := range buildInfo.ReadySCs {
		isLeader := scInfo.
			Address.
			Attributes.
			Value("is_leader").(bool)

		if isLeader {
			p.leader = sc
			continue
		}
		followers = append(followers, sc)
	}
	p.followers = followers
	return p
}

var _ balancer.Picker = (*Picker)(nil)

func (p *Picker) Pick(info balancer.PickInfo) (
	balancer.PickResult, error,
) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result balancer.PickResult
	if strings.Contains(info.FullMethodName, "SaveMessage") ||
		 strings.Contains(info.FullMethodName, "DeleteMessages") ||
		 strings.Contains(info.FullMethodName, "DeleteUserMessages") ||
		 strings.Contains(info.FullMethodName, "PublishMessages") ||
		 strings.Contains(info.FullMethodName, "PrivateMessages") ||
		 strings.Contains(info.FullMethodName, "UpdateMessage") ||
		len(p.followers) == 0 {
			result.SubConn = p.leader
	} else if strings.Contains(info.FullMethodName, "ReadMessages") ||
						strings.Contains(info.FullMethodName, "ReadThreadMessages") ||
						strings.Contains(info.FullMethodName, "ReadMessage") ||
						strings.Contains(info.FullMethodName, "ReadPath") ||
						strings.Contains(info.FullMethodName, "ReadMessagesAround") ||
						strings.Contains(info.FullMethodName, "CountMessages") {
		result.SubConn = p.nextFollower()
	}
	if result.SubConn == nil {
		return result, balancer.ErrNoSubConnAvailable
	}
	return result, nil
}

func (p *Picker) nextFollower() balancer.SubConn {
	cur := atomic.AddUint64(&p.current, uint64(1))
	len := uint64(len(p.followers))
	idx := int(cur % len)
	return p.followers[idx]
}

func init() {
	balancer.Register(
		base.NewBalancerBuilder(Name, &Picker{}, base.Config{}),
	)
}