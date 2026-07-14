package config

type Config struct {
	NodeName            string   `json:"node_name"`
	RpcAddr             string   `json:"rpc_addr"`
	SerfAddr            string   `json:"serf_addr"`
	RaftServers         []string `json:"raft_servers"`
	SerfJoinAddrs       []string `json:"serf_join_addrs"`
	RaftLogLevel        string   `json:"raft_log_level"`
	LogLevel            string   `json:"log_level"`

	RaftBootstrap bool   `json:"raft_bootstrap"`
	PGConn        string `json:"pg_conn"`
	// TODO ShutdownTimeout string     `json:"shutdown_timeout"`
	GooseTableName string `json:"goose_table_name"`
	DataPath       string `json:"data_path"`
}
