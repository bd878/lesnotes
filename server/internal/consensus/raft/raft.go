package raft

import (
	"os"
	"time"
	"errors"
	"context"
	"path/filepath"

	raftboltdb "github.com/hashicorp/raft-boltdb"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
)

type Config struct {
	Bootstrap            bool
	HeartbeatTimeout     time.Duration
	ElectionTimeout      time.Duration
	LeaderLeaseTimeout   time.Duration
	CommitTimeout        time.Duration
	NodeName             string
	RaftLogLevel         string
	DataDir              string
	Servers              []string
	RetainSnapshots      int
	MaxConnectionsPool   int
	NetworkTimeout       time.Duration
}

type Distributed struct {
	log            *logger.Logger
	conf           Config
	raft           *raft.Raft
	snapshotStore  raft.SnapshotStore
	streamLayer    *StreamLayer
}

func New(conf Config, streamLayer *StreamLayer, fsm raft.FSM, log *logger.Logger) (*Distributed, error) {
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
		log:         log,
		conf:        conf,
		streamLayer: streamLayer,
	}

	raftLogLevel := hclog.Error.String()
	switch conf.RaftLogLevel {
	case "debug":
		raftLogLevel = hclog.Debug.String()
	case "error":
		raftLogLevel = hclog.Error.String()
	case "info":
		raftLogLevel = hclog.Info.String()
	default:
		raftLogLevel = hclog.Info.String()
	}

	raftPath := filepath.Join(conf.DataDir, "raft")
	if err := os.MkdirAll(raftPath, 0755); err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "log"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftPath, "stable"))
	if err != nil {
		return nil, err
	}

	m.snapshotStore, err = raft.NewFileSnapshotStore(filepath.Join(raftPath, "snapshot"), conf.RetainSnapshots, os.Stderr)
	if err != nil {
		return nil, err
	}

	transport := raft.NewNetworkTransport(streamLayer, conf.MaxConnectionsPool, conf.NetworkTimeout, os.Stderr)

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(conf.NodeName)
	config.LogLevel = raftLogLevel
	if conf.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = conf.HeartbeatTimeout
	}
	if conf.ElectionTimeout != 0 {
		config.ElectionTimeout = conf.ElectionTimeout
	}
	if conf.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = conf.LeaderLeaseTimeout
	}
	if conf.CommitTimeout != 0 {
		config.CommitTimeout = conf.CommitTimeout
	}
	if conf.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = conf.LeaderLeaseTimeout
	}

	m.raft, err = raft.NewRaft(config, fsm, logStore, stableStore, m.snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	hasState, err := raft.HasExistingState(logStore, stableStore, m.snapshotStore)
	if err != nil {
		return nil, err
	}

	if conf.Bootstrap && !hasState {
		servers := []raft.Server{{
			ID:      raft.ServerID(conf.NodeName),
			Address: transport.LocalAddr(),
		}}

		for _, addr := range conf.Servers {
			servers = append(servers, raft.Server{
				ID:      raft.ServerID(addr),
				Address: raft.ServerAddress(addr),
			})
		}

		configuration := raft.Configuration{
			Servers: servers,
		}
		err = m.raft.BootstrapCluster(configuration).Error()
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Distributed) WaitForLeader(timeout time.Duration) error {
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- timeoutc:
			m.log.Error("no leader, timeout")
			return nil
		case <-ticker.C:
			if lead, _ := m.raft.LeaderWithID(); lead != "" {
				return nil
			}
		}
	}
}

func (m *Distributed) Apply(cmd []byte, timeout time.Duration) (err error) {
	future := m.raft.Apply(cmd, timeout)
	if future.Error() != nil {
		return future.Error()
	}

	res := future.Response()
	if err, ok := res.(error); ok {
		return err
	}

	return nil
}

func (m *Distributed) GetServers(_ context.Context) ([](*api.Server), error) {
	future := m.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		m.log.Error("message", "failed to get servers configuration")
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

	m.log.Infow("remove from cluster server", "id", id)
	removeFuture := m.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (m *Distributed) Snapshot() error {
	m.log.Debugln("snapshot this machine")

	snapshotFuture := m.raft.Snapshot()
	return snapshotFuture.Error()
}

func (m *Distributed) Restore() error {
	leaderFuture := m.raft.VerifyLeader()
	if err := leaderFuture.Error(); err != nil {
		return errors.New("cannot restore from snapshot: not a leader")
	}

	m.log.Debugln("restoring from last snapshot")
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
	m.log.Infow("my state", "addr", m.streamLayer.Addr(), "state", state.String())
	return nil
}

func (m *Distributed) NodeName() string {
	return m.conf.NodeName
}

func (m *Distributed) isLeader() bool {
	return m.raft.State() == raft.Leader
}