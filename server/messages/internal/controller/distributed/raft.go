package messages

import (
  "io"
  "os"
  "fmt"
  "net"
  "time"
  "bytes"
  "errors"
  "encoding/json"
  "context"
  "path/filepath"

  raftboltdb "github.com/hashicorp/raft-boltdb"
  "github.com/hashicorp/raft"

  usermodel "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/messages/internal/repository"

  "github.com/bd878/gallery/server/api"
)

var ErrMsgExist = errors.New("message exists")

type Repository interface {
  Put(context.Context, *model.Message) (model.MessageId, error)
  Get(context.Context, usermodel.UserId) ([]*model.Message, error)
  FindByIndexTerm(context.Context, uint64, uint64) (*model.Message, error)
  PutBatch(context.Context, [](*model.Message)) error
  GetBatch(context.Context) ([]*model.Message, error)
  GetOne(context.Context, usermodel.UserId, model.MessageId) (*model.Message, error)
  Truncate(context.Context) error
}

type DistributedMessages struct {
  config   Config
  raft    *raft.Raft
  repo     Repository
}

func New(repo Repository, config Config) (
  *DistributedMessages,
  error,
) {
  m := &DistributedMessages{
    repo: repo,
    config: config,
  }
  if err := m.setupRaft(); err != nil {
    return nil, err
  }
  return m, nil
}

func (m *DistributedMessages) setupRaft() error {
  fsm := &fsm{repo: m.repo}

  raftPath := filepath.Join(m.config.DataDir, "raft")
  if err := os.MkdirAll(raftPath, 0755); err != nil {
    return err
  }

  logStore, err := raftboltdb.NewBoltStore(
    filepath.Join(raftPath, "log"),
  )
  if err != nil {
    return err
  }
  stableStore, err := raftboltdb.NewBoltStore(
    filepath.Join(raftPath, "stable"),
  )
  if err != nil {
    return err
  }
  retain := 1
  snapshotStore, err := raft.NewFileSnapshotStore(
    filepath.Join(raftPath, "raft"),
    retain,
    nil,
  )
  if err != nil {
    return err
  }

  maxPool := 5
  timeout := 10*time.Second
  transport := raft.NewNetworkTransport(
    m.config.StreamLayer,
    maxPool,
    timeout,
    os.Stderr,
  )

  config := raft.DefaultConfig()
  config.LocalID = m.config.Raft.LocalID
  config.LogLevel = m.config.Raft.LogLevel
  if m.config.Raft.HeartbeatTimeout != 0 {
    config.HeartbeatTimeout = m.config.Raft.HeartbeatTimeout
  }
  if m.config.Raft.ElectionTimeout != 0 {
    config.ElectionTimeout = m.config.Raft.ElectionTimeout
  }
  if m.config.Raft.LeaderLeaseTimeout != 0 {
    config.LeaderLeaseTimeout = m.config.Raft.LeaderLeaseTimeout
  }
  if m.config.Raft.CommitTimeout != 0 {
    config.CommitTimeout = m.config.Raft.CommitTimeout
  }
  if m.config.Raft.LeaderLeaseTimeout != 0 {
    config.LeaderLeaseTimeout = m.config.Raft.LeaderLeaseTimeout
  }

  m.raft, err = raft.NewRaft(
    config,
    fsm,
    logStore,
    stableStore,
    snapshotStore,
    transport,
  )
  if err != nil {
    return err
  }

  var hasState bool
  hasState, err = raft.HasExistingState(
    logStore,
    stableStore,
    snapshotStore,
  )
  if err != nil {
    return err
  }
  if m.config.Bootstrap && !hasState {
    servers := []raft.Server{{
      ID: m.config.Raft.LocalID,
      Address: transport.LocalAddr(),
    }}

    for _, addr := range m.config.Servers {
      servers = append(servers, raft.Server{
        ID: raft.ServerID(addr),
        Address: raft.ServerAddress(addr),
      })
    }

    configuration := raft.Configuration{
      Servers: servers,
    }
    err = m.raft.BootstrapCluster(configuration).Error()
  }
  return err
}

func (m *DistributedMessages) SaveMessage(ctx context.Context, msg *model.Message) (resMsg *model.Message, err error) {
  msg.CreateTime = time.Now().String()
  if resMsg, err = m.apply(ctx, msg); err != nil {
    return nil, err
  }
  return resMsg, nil
}

func (m *DistributedMessages) apply(ctx context.Context, msg *model.Message) (*model.Message, error) {
  b, err := json.Marshal(msg)
  if err != nil {
    return nil, err
  }

  timeout := 10*time.Second
  future := m.raft.Apply(b, timeout)
  if future.Error() != nil {
    return nil, future.Error()
  }

  res := future.Response()
  switch val := res.(type) {
  case error:
    return nil, val
  case model.Message:
    return &val, nil
  default:
    return nil, errors.New("fsm.apply returns undefined result")
  }
}

func (m *DistributedMessages) ReadUserMessages(ctx context.Context, userId usermodel.UserId) (
  []*model.Message,
  error,
) {
  return m.repo.Get(ctx, userId)
}

func (m *DistributedMessages) ReadOneMessage(ctx context.Context, userId usermodel.UserId, id model.MessageId) (
  *model.Message,
  error,
) {
  return m.repo.GetOne(ctx, userId, id)
}

func (m *DistributedMessages) WaitForLeader(timeout time.Duration) error {
  timeoutc := time.After(timeout)
  ticker := time.NewTicker(time.Second)
  defer ticker.Stop()
  for {
    select {
    case <- timeoutc:
      return fmt.Errorf("no leader, timeout")
    case <-ticker.C:
      if lead, _ := m.raft.LeaderWithID(); lead != "" {
        return nil
      }
    }
  }
} 

func (m *DistributedMessages) GetServers(_ context.Context) ([](*api.Server), error) {
  future := m.raft.GetConfiguration()
  if err := future.Error(); err != nil {
    return nil, err
  }
  var servers []*api.Server
  leaderAddr, _ := m.raft.LeaderWithID()
  for _, server := range future.Configuration().Servers {
    servers = append(servers, &api.Server{
      Id: string(server.ID),
      RaftAddr: string(server.Address),
      IsLeader: leaderAddr == server.Address,
    })
  }
  return servers, nil
}

func (m *DistributedMessages) Join(id, addr string) error {
  leaderFuture := m.raft.VerifyLeader()
  if err := leaderFuture.Error(); err != nil {
    return errors.New("cannot join node to cluster: not a leader")
  }

  configFuture := m.raft.GetConfiguration()
  if err := configFuture.Error(); err != nil {
    return err
  }

  serverID := raft.ServerID(id)
  serverAddr := raft.ServerAddress(addr)

  for _, srv := range configFuture.Configuration().Servers {
    if srv.ID == serverID || srv.Address == serverAddr {
      if srv.ID == serverID && srv.Address == serverAddr {
        return nil
      }

      removeFuture := m.raft.RemoveServer(serverID, 0, 0)
      if err := removeFuture.Error(); err != nil {
        return err
      }
    }
  }

  addFuture := m.raft.AddVoter(serverID, serverAddr, 0, 0)
  if err := addFuture.Error(); err != nil {
    return err
  }

  return nil
}

func (m *DistributedMessages) Leave(id string) error {
  leaderFuture := m.raft.VerifyLeader()
  if err := leaderFuture.Error(); err != nil {
    return errors.New("cannot remove node from cluster: not a leader")
  }

  removeFuture := m.raft.RemoveServer(raft.ServerID(id), 0, 0)
  return removeFuture.Error()
}

func (m *DistributedMessages) PrintLeader() error {
  addr, id := m.raft.LeaderWithID()
  fmt.Println("=== LEADER ===")
  if m.config.Raft.LocalID == raft.ServerID(id) {
    fmt.Println("i am the leader")
  }
  fmt.Printf("Addr: %v\n", addr)
  fmt.Printf("Id: %v\n", id)
  fmt.Println()
  return nil
}

func (m *DistributedMessages) PrintMyAddr() error {
  _, id := m.raft.LeaderWithID()
  addr := m.config.StreamLayer.Addr()
  fmt.Println("=== ME ===")
  if m.config.Raft.LocalID == raft.ServerID(id) {
    fmt.Println("i am the leader")
  }
  fmt.Printf("Address: %v\n", addr)
  fmt.Printf("ID: %v\n", m.config.Raft.LocalID)
  fmt.Println()
  return nil
}

func (m *DistributedMessages) PrintConfig() error {
  future := m.raft.GetConfiguration()
  err := future.Error()
  if err != nil {
    return err
  }

  fmt.Println("=== SERVERS ===")
  conf := future.Configuration()
  for i, serv := range conf.Servers {
    fmt.Printf("# %d:\n", i)
    fmt.Printf("Suffrage: %d\n", serv.Suffrage)
    fmt.Printf("Id: %s\n", serv.ID)
    fmt.Printf("Address: %s\n", serv.Address)
    fmt.Println()
  }

  return nil
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
  repo Repository
}

/**
 * Returns empty interface. It is either an error,
 * or new msg with unique id, saved in repo.
 * 
 * Apply replicates log state from the bottom up.
 * Leader makes Apply on start.
 */
func (f *fsm) Apply(record *raft.Log) interface{} {
  var msg *model.Message
  var err error

  msg, err = f.repo.FindByIndexTerm(context.Background(), record.Index, record.Term)
  if err != nil {
    /* not found is expected behaviour */
    if !errors.Is(err, repository.ErrNotFound) {
      return err
    }
  }
  if msg != nil {
    return ErrMsgExist
  }

  buf := record.Data
  err = json.Unmarshal(buf, &msg)
  if err != nil {
    return err
  }
  msg.LogIndex = record.Index
  msg.LogTerm = record.Term

  msg.Id, err = f.repo.Put(context.Background(), msg)
  if err != nil {
    return err
  }

  return *msg
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
  return &snapshot{repo: f.repo}, nil
}

// TODO: restore will reapply same messages with ids,
// check whether msg with id exists
func (f *fsm) Restore(r io.ReadCloser) error {
  var buf *bytes.Buffer
  var msgs []model.Message

  _, err := io.Copy(buf, r)
  if err == io.EOF {
    return err
  } else if err != nil {
    return err
  }
  err = json.Unmarshal(buf.Bytes(), &msgs)
  if err != nil {
    return err
  }

  ctx := context.Background()
  err = f.repo.Truncate(ctx)
  if err != nil {
    return err
  }
  for _, msg := range msgs {
    _, err := f.repo.Put(ctx, &msg)
    if err != nil {
      return err
    }
  }
  return nil
}

type snapshot struct {
  repo Repository
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
  msgs, err := s.repo.GetBatch(context.Background())
  if err != nil {
    return err
  }

  b, err := json.Marshal(msgs)
  if err != nil {
    return err
  }
  if _, err := io.Copy(sink, bytes.NewReader(b)); err != nil {
    _ = sink.Cancel()
    return err
  }
  return sink.Close()
}

func (s *snapshot) Release() {}

type StreamLayer struct {
  ln net.Listener
}

func NewStreamLayer(ln net.Listener) *StreamLayer {
  return &StreamLayer{ln: ln}
}

func (s *StreamLayer) Dial(
  addr raft.ServerAddress,
  timeout time.Duration,
) (net.Conn, error) {
  dialer := &net.Dialer{Timeout: timeout}
  return dialer.Dial("tcp", string(addr))
}

func (s *StreamLayer) Accept() (net.Conn, error) {
  return s.ln.Accept()
}

func (s *StreamLayer) Close() error {
  return s.ln.Close()
}

func (s *StreamLayer) Addr() net.Addr {
  return s.ln.Addr()
}