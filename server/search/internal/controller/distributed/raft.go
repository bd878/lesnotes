package distributed

import (
	"os"
	"time"
	"errors"
	"context"
	"path/filepath"

	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/raft"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
)

var ErrMsgExist = errors.New("message exists")

type Config struct {
	Raft                 raft.Config
	StreamLayer          *StreamLayer
	Bootstrap            bool
	DataDir              string
	Servers              []string
	RetainSnapshots      int
	MaxConnectionsPool   int
	NetworkTimeout       time.Duration
}

type Distributed struct {
	conf            Config
	raft            *raft.Raft
	repo            Repository
	snapshotStore   raft.SnapshotStore
}

func New(conf Config, repo Repository) (*Distributed, error) {
	if conf.RetainSnapshots == 0 {
		conf.RetainSnapshots = 1
	}

	if conf.MaxConnectionsPool == 0 {
		conf.MaxConnectionsPool = 5
	}

	if conf.NetworkTimeout == 0 {
		conf.NetworkTimeout = 10 * time.Second
	}

	m := &Distributed{
		repo:      repo,
		conf:      conf,
	}
	if err := m.setupRaft(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Distributed) setupRaft() error {
	fsm := &fsm{repo: m.repo}

	raftPath := filepath.Join(m.conf.DataDir, "raft")
	if err := os.MkdirAll(raftPath, 0755); err != nil {
		return err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "log"))
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "stable"))
	if err != nil {
		return err
	}

	m.snapshotStore, err = raft.NewFileSnapshotStore(filepath.Join(raftPath, "snapshot"), m.conf.RetainSnapshots, os.Stderr)
	if err != nil {
		return err
	}

	transport := raft.NewNetworkTransport(m.conf.StreamLayer, m.conf.MaxConnectionsPool, m.conf.NetworkTimeout, os.Stderr)

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

	m.raft, err = raft.NewRaft(config, fsm, logStore,
		stableStore, m.snapshotStore, transport)
	if err != nil {
		return err
	}

	var hasState bool
	hasState, err = raft.HasExistingState(logStore, stableStore, m.snapshotStore)
	if err != nil {
		return err
	}
	if m.conf.Bootstrap && !hasState {
		servers := []raft.Server{{
			ID:      m.conf.Raft.LocalID,
			Address: transport.LocalAddr(),
		}}

		for _, addr := range m.conf.Servers {
			servers = append(servers, raft.Server{
				ID:      raft.ServerID(addr),
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

func (m *Distributed) WaitForLeader(timeout time.Duration) error {
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

func (m *Distributed) GetServers(_ context.Context) ([](*api.Server), error) {
	future := m.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		logger.Error("message", "failed to get servers configuration")
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

func (m *Distributed) Join(id, addr string) error {
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

func (m *Distributed) Leave(id string) error {
	leaderFuture := m.raft.VerifyLeader()
	if err := leaderFuture.Error(); err != nil {
		return errors.New("cannot remove node from cluster: not a leader")
	}

	logger.Info("remove from cluster serve with id", id)
	removeFuture := m.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (m *Distributed) Snapshot() error {
	logger.Debugln("snapshot this machine")

	snapshotFuture := m.raft.Snapshot()
	return snapshotFuture.Error()
}

func (m *Distributed) Restore() error {
	leaderFuture := m.raft.VerifyLeader()
	if err := leaderFuture.Error(); err != nil {
		return errors.New("cannot restore from snapshot: not a leader")
	}

	logger.Debugln("restoring from last snapshot")
	list, err := m.snapshotStore.List()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return errors.New("cannot restore from snapshot: no snapshots")
	}

	snapshot, reader, err := m.snapshotStore.Open(list[0].ID)
	defer reader.Close()
	if err != nil {
		return err
	}

	return m.raft.Restore(snapshot, reader, 20 * time.Second)
}

func (m *Distributed) ShowLeader() error {
	state := m.raft.State()
	logger.Infow("my state", "addr", m.conf.StreamLayer.Addr(), "state", state.String())
	return nil
}