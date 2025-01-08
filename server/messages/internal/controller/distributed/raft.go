package messages

import (
  "io"
  "os"
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
  "github.com/bd878/gallery/server/logger"
)

var ErrMsgExist = errors.New("message exists")

type Config struct {
  Raft         raft.Config
  StreamLayer *StreamLayer
  Bootstrap    bool
  DataDir      string
  Servers      []string
}

type Repository interface {
  Put(ctx context.Context, log *logger.Logger, params *model.PutParams) (model.MessageId, error)
  Get(ctx context.Context, log *logger.Logger, params *model.GetParams) (*model.MessagesList, error)
  FindByIndexTerm(ctx context.Context, log *logger.Logger, params *model.FindByIndexParams) (*model.Message, error)
  GetBatch(ctx context.Context, log *logger.Logger) ([]*model.Message, error)
  GetOne(ctx context.Context, log *logger.Logger, params *model.GetOneParams) (*model.Message, error)
  Truncate(ctx context.Context, log *logger.Logger) error
}

type DistributedMessages struct {
  conf     Config
  raft    *raft.Raft
  repo     Repository
}

func New(repo Repository, cfg Config) (*DistributedMessages, error) {
  m := &DistributedMessages{
    repo: repo,
    conf: cfg,
  }
  if err := m.setupRaft(logger.Default()); err != nil {
    return nil, err
  }
  return m, nil
}

func (m *DistributedMessages) setupRaft(log *logger.Logger) error {
  fsm := &fsm{repo: m.repo}

  raftPath := filepath.Join(m.conf.DataDir, "raft")
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
    m.conf.StreamLayer,
    maxPool,
    timeout,
    os.Stderr,
  )

  config := raft.DefaultConfig()
  config.LocalID = m.conf.Raft.LocalID
  config.LogLevel = m.conf.Raft.LogLevel
  if m.conf.Raft.HeartbeatTimeout != 0 {
    config.HeartbeatTimeout = m.conf.Raft.HeartbeatTimeout
  }
  if m.conf.Raft.ElectionTimeout != 0 {
    config.ElectionTimeout = m.conf.Raft.ElectionTimeout
  }
  if m.conf.Raft.LeaderLeaseTimeout != 0 {
    config.LeaderLeaseTimeout = m.conf.Raft.LeaderLeaseTimeout
  }
  if m.conf.Raft.CommitTimeout != 0 {
    config.CommitTimeout = m.conf.Raft.CommitTimeout
  }
  if m.conf.Raft.LeaderLeaseTimeout != 0 {
    config.LeaderLeaseTimeout = m.conf.Raft.LeaderLeaseTimeout
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
  if m.conf.Bootstrap && !hasState {
    servers := []raft.Server{{
      ID: m.conf.Raft.LocalID,
      Address: transport.LocalAddr(),
    }}

    for _, addr := range m.conf.Servers {
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

func (m *DistributedMessages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (
  resMsg *model.Message, err error,
) {
  params.Message.CreateTime = time.Now().String()
  if resMsg, err = m.apply(ctx, params.Message); err != nil {
    log.Error("message", "raft failed to apply save message")
    return nil, err
  }
  return resMsg, nil
}

func (m *DistributedMessages) UpdateMessage(_ context.Context, _ *logger.Logger, _ *model.UpdateMessageParams) (
  *model.Message, error,
) {
  /* not implemented */
  return nil, nil
}

func (m *DistributedMessages) apply(ctx context.Context, msg *model.Message) (
  *model.Message, error,
) {
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

func (m *DistributedMessages) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.MessagesList, error,
) {
  return m.repo.Get(
    ctx,
    log,
    &model.GetParams{
      UserId:    params.UserId,
      Limit:     params.Limit,
      Offset:    params.Offset,
      Ascending: params.Ascending,
    },
  )
}

func (m *DistributedMessages) ReadOneMessage(ctx context.Context, userId usermodel.UserId, id model.MessageId) (
  *model.Message, error,
) {
  return m.repo.GetOne(ctx, logger.Default(), &model.GetOneParams{userId, id})
}

func (m *DistributedMessages) WaitForLeader(timeout time.Duration) error {
  timeoutc := time.After(timeout)
  ticker := time.NewTicker(time.Second)
  defer ticker.Stop()
  for {
    select {
    case <- timeoutc:
      logger.Error("no leader, timeout")
      return nil
    case <-ticker.C:
      if lead, _ := m.raft.LeaderWithID(); lead != "" {
        return nil
      }
    }
  }
}

func (m *DistributedMessages) GetServers(_ context.Context, log *logger.Logger) ([](*api.Server), error) {
  future := m.raft.GetConfiguration()
  if err := future.Error(); err != nil {
    log.Error("message", "failed to get servers configuration")
    return nil, err
  }
  var servers []*api.Server
  _, id := m.raft.LeaderWithID()
  for _, server := range future.Configuration().Servers {
    servers = append(servers, &api.Server{
      Id: string(server.ID),
      RaftAddr: string(server.Address),
      IsLeader: raft.ServerID(id) == server.ID,
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

  logger.Info("remove from cluster serve with id", id)
  removeFuture := m.raft.RemoveServer(raft.ServerID(id), 0, 0)
  return removeFuture.Error()
}

func (m *DistributedMessages) PrintLeader() error {
  addr, id := m.raft.LeaderWithID()
  logger.Infoln("=== LEADER ===")
  if m.conf.Raft.LocalID == raft.ServerID(id) {
    logger.Infoln("i am the leader")
  }
  logger.Info("Addr: %v\n", addr)
  logger.Info("Id: %v\n", id)
  logger.Infoln()
  return nil
}

func (m *DistributedMessages) PrintMyAddr() error {
  _, id := m.raft.LeaderWithID()
  addr := m.conf.StreamLayer.Addr()
  logger.Infoln("=== ME ===")
  if m.conf.Raft.LocalID == raft.ServerID(id) {
    logger.Infoln("i am the leader")
  }
  logger.Infoln("Address: %v\n", addr)
  logger.Infoln("ID: %v\n", m.conf.Raft.LocalID)
  logger.Infoln()
  return nil
}

func (m *DistributedMessages) PrintConfig() error {
  future := m.raft.GetConfiguration()
  err := future.Error()
  if err != nil {
    return err
  }

  logger.Infoln("=== SERVERS ===")
  conf := future.Configuration()
  for i, serv := range conf.Servers {
    logger.Info("# %d:\n", i)
    logger.Info("Suffrage: %d\n", serv.Suffrage)
    logger.Info("Id: %s\n", serv.ID)
    logger.Info("Address: %s\n", serv.Address)
    logger.Infoln()
  }

  return nil
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
  repo Repository
}

type RequestType uint8

const (
  AppendRequestType RequestType = 0
  UpdateRequestType
)

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

  msg, err = f.repo.FindByIndexTerm(context.Background(), logger.Default(), &model.FindByIndexParams{
    LogIndex: record.Index,
    LogTerm:  record.Term,
  })
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

  msg.Id, err = f.repo.Put(context.Background(), logger.Default(), &model.PutParams{Message: msg})
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
  err = f.repo.Truncate(ctx, logger.Default())
  if err != nil {
    return err
  }
  for _, msg := range msgs {
    _, err := f.repo.Put(ctx, logger.Default(), &model.PutParams{Message: &msg})
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
  msgs, err := s.repo.GetBatch(context.Background(), logger.Default())
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

/**
 * Leading byte in rpc, signifies, that it is a call
 * to raft node (not grpc)
 */
const RaftRPC = 1

func (s *StreamLayer) Dial(
  addr raft.ServerAddress,
  timeout time.Duration,
) (net.Conn, error) {
  dialer := &net.Dialer{Timeout: timeout}
  conn, err := dialer.Dial("tcp", string(addr))
  if err != nil {
    return nil, err
  }

  _, err = conn.Write([]byte{byte(RaftRPC)})
  if err != nil {
    return nil, err
  }

  return conn, err
}

func (s *StreamLayer) Accept() (net.Conn, error) {
  conn, err := s.ln.Accept()
  if err != nil {
    return nil, err
  }

  b := make([]byte, 1)
  if _, err := conn.Read(b); err != nil {
    return nil, err
  }
  if bytes.Compare(b, []byte{byte(RaftRPC)}) != 0 {
    return nil, errors.New("not a raft rpc")
  }

  return conn, nil
}

func (s *StreamLayer) Close() error {
  return s.ln.Close()
}

func (s *StreamLayer) Addr() net.Addr {
  return s.ln.Addr()
}