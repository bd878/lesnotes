package config

type Config struct {
	NodeName             string      `json:"node_name"`
	HttpAddr             string      `json:"http_addr"`
	UsersServiceAddr     string      `json:"users_service_addr"`
	FilesServiceAddr     string      `json:"files_service_addr"`
	SessionsServiceAddr  string      `json:"sessions_service_addr"`
	RpcAddr              string      `json:"rpc_addr"`
	SerfAddr             string      `json:"serf_addr"`
	RaftServers          []string    `json:"raft_servers"`
	SerfJoinAddrs        []string    `json:"serf_join_addrs"`
	RaftLogLevel         string      `json:"raft_log_level"`
	LogLevel             string      `json:"log_level"`

	RaftBootstrap        bool        `json:"raft_bootstrap"`
	LogPath              string      `json:"log_path"`
	PGConn               string      `json:"pg_conn"`
	// TODO ShutdownTimeout string     `json:"shutdown_timeout"`
	DataPath             string      `json:"data_path"`
}
