package loadbalance_test

import (
  "testing"

  "github.com/stretchr/testify/require"
  "google.golang.org/grpc/resolver"
  "google.golang.org/grpc/balancer"
  "google.golang.org/grpc/balancer/base"
  "google.golang.org/grpc/attributes"

  "github.com/bd878/gallery/server/messages/internal/loadbalance"
)

func TestPickerSaveToLeader(t *testing.T) {
  picker, subConns := setupTest()
  info := balancer.PickInfo{
    FullMethodName: "/messages.v1.MessagesService/SaveMessage",
  }
  for i := 0; i < 5; i++ {
    gotPick, err := picker.Pick(info)
    require.NoError(t, err)
    require.Equal(t, subConns[0], gotPick.SubConn)
  }
}

func setupTest() (*loadbalance.Picker, []*subConn) {
  readySCs := make(map[balancer.SubConn]base.SubConnInfo)
  var subConns []*subConn
  for i := 0; i < 3; i++ {
    addr := resolver.Address{
      Attributes: attributes.New("is_leader", i == 0),
    }

    sc := &subConn{addrs: []resolver.Address{addr}}
    readySCs[sc] = base.SubConnInfo{Address:addr}
    subConns = append(subConns, sc)
  }

  picker := &loadbalance.Picker{}
  picker.Build(base.PickerBuildInfo{
    ReadySCs: readySCs,
  })
  return picker, subConns
}

type subConn struct {
  addrs []resolver.Address
}

var _ balancer.SubConn = (*subConn)(nil)

func (c *subConn) Connect() {}

func (c *subConn) Shutdown() {}

func (c *subConn) UpdateAddresses(addrs []resolver.Address) {
  c.addrs = addrs
}

func (c *subConn) GetOrBuildProducer(b balancer.ProducerBuilder) (
  balancer.Producer, func(),
) {
  return b.Build(&struct{}{})
}