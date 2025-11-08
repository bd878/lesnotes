package config

type Config struct {
	RpcAddr           string          `json:"rpc_addr"`
	NodeName          string          `json:"node_name"`

	SerfAddr               string          `json:"serf_addr"`
	RaftServers            []string        `json:"raft_servers"`
	SerfJoinAddrs          []string        `json:"serf_join_addrs"`
	RaftLogLevel           string          `json:"raft_log_level"`
	LogLevel               string          `json:"log_level"`
	NatsAddr               string          `json:"nats_addr"`

	RaftBootstrap          bool            `json:"raft_bootstrap"`
	DataPath               string          `json:"data_path"`
	PGConn                 string          `json:"pg_conn"`
	TableName              string          `json:"table_name"`
}
