package messages_test

import (
  "testing"
  "context"
  "reflect"
  "time"
  "net"
  "fmt"

  "github.com/hashicorp/raft"
  "github.com/stretchr/testify/require"

  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/user/pkg/model"
  memory "github.com/bd878/gallery/server/messages/internal/repository/memory"
  distributed "github.com/bd878/gallery/server/messages/internal/controller/distributed"
)

var ports []int = []int{8081, 8082, 8083, 8084, 8085}

func TestDistributed(t *testing.T) {
  nodeCount := len(ports)

  var logs []*distributed.DistributedMessages
  for i := 0; i < nodeCount; i++ {
    dataDir := t.TempDir()
    repo := memory.New()

    ln, err := net.Listen("tcp",
      fmt.Sprintf("127.0.0.1:%d", ports[i]),
    )
    require.NoError(t, err)

    config := distributed.Config{}
    config.StreamLayer = distributed.NewStreamLayer(ln)
    config.Raft.LocalID = raft.ServerID(fmt.Sprintf("%d", i))
    config.DataDir = dataDir
    config.Raft.HeartbeatTimeout = 50 * time.Millisecond
    config.Raft.ElectionTimeout = 50 * time.Millisecond
    config.Raft.LeaderLeaseTimeout = 20 * time.Millisecond
    config.Raft.CommitTimeout = 5 * time.Millisecond

    if i == 0 {
      config.Bootstrap = true
    }

    m, err := distributed.New(repo, config)
    require.NoError(t, err)

    if i != 0 {
      err = logs[0].Join(
        fmt.Sprintf("%d", i), ln.Addr().String(),
      )
      require.NoError(t, err)
    } else {
      err = m.WaitForLeader(3 * time.Second)
      require.NoError(t, err)
    }

    logs = append(logs, m)
  }

  messages := []*model.Message{
    {Id: 0, UserId: 1, Value: "first", File: "file1_1.pdf"},
    {Id: 1, UserId: 2, Value: "second", File: "file2_1.pdf"},
    {Id: 2, UserId: 1, Value: "third", File: "file1_2.pdf"},
  }

  for _, msg := range messages {
    err := logs[0].SaveMessage(context.Background(), msg)
    require.NoError(t, err)
    require.Eventually(t, func() bool {
      for j := 0; j < nodeCount; j++ {
        got, err := logs[j].ReadOneMessage(
          context.Background(),
          usermodel.UserId(msg.UserId),
          msg.Id,
        )
        if err != nil {
          return false
        }

        if !reflect.DeepEqual(got.Value, msg.Value) {
          return false
        }
      }
      return true
    }, 500*time.Millisecond, 50*time.Millisecond)
  }

  servers, err := logs[0].GetServers()
  require.NoError(t, err)
  require.Equal(t, nodeCount, len(servers))
  require.True(t, servers[0].IsLeader)
  require.False(t, servers[1].IsLeader)
  require.False(t, servers[2].IsLeader)

  err = logs[0].Leave("1")
  require.NoError(t, err)

  time.Sleep(50 *time.Millisecond)

  servers, err = logs[0].GetServers()
  require.NoError(t, err)
  require.Equal(t, nodeCount-1, len(servers))

  err = logs[0].SaveMessage(context.Background(), &model.Message{
    Id: 3,
    UserId: 1,
    Value: "third",
  })
  require.NoError(t, err)

  time.Sleep(50 * time.Millisecond)
  message, err := logs[2].ReadOneMessage(context.Background(), usermodel.UserId(1), 3)
  require.NoError(t, err)
  require.Equal(t, "third", message.Value)
}