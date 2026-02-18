package config

type Config struct {
	NodeName             string      `json:"node_name"`
	RpcAddr              string      `json:"rpc_addr"`
	HttpAddr             string      `json:"http_addr"`
	SessionsServiceAddr  string      `json:"sessions_service_addr"`
	UsersServiceAddr     string      `json:"users_service_addr"`

	SerfAddr             string      `json:"serf_addr"`
	RaftServers          []string    `json:"raft_servers"`
	SerfJoinAddrs        []string    `json:"serf_join_addrs"`
	RaftLogLevel         string      `json:"raft_log_level"`
	LogLevel             string      `json:"log_level"`
	NatsAddr             string      `json:"nats_addr"`

	RaftBootstrap        bool        `json:"raft_bootstrap"`
	DataPath             string      `json:"data_path"`
	PGConn               string      `json:"pg_conn"`
	TableName            string      `json:"table_name"`
	GooseTableName       string      `json:"goose_table_name"`
}
