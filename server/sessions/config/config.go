package config

type Config struct {
	RpcAddr           string          `json:"rpc_addr"`
	NodeName          string          `json:"node_name"`
	LogPath           string          `json:"log_path"`
	PGConn            string          `json:"pg_conn"`
}
