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
	Raft         raft.Config
	StreamLayer *StreamLayer
	Bootstrap    bool
	DataDir      string
	Servers      []string
	DBPath       string
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

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "log"))
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "stable"))
	if err != nil {
		return err
	}
	// TODO: rewrite on SqliteSnapshotStore from sqlite_snapshot branch
	snapshotStore := raft.NewDiscardSnapshotStore()

	maxPool := 5
	timeout := 10*time.Second
	transport := raft.NewNetworkTransport(m.conf.StreamLayer, maxPool, timeout,
		os.Stderr)

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
		stableStore, snapshotStore, transport)
	if err != nil {
		return err
	}

	var hasState bool
	hasState, err = raft.HasExistingState(logStore, stableStore, snapshotStore)
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

func (m *DistributedMessages) GetServers(_ context.Context) ([](*api.Server), error) {
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
