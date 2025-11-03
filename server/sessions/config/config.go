package config

type Config struct {
	RpcAddr           string          `json:"rpc_addr"`
	NodeName          string          `json:"node_name"`
	PGConn            string          `json:"pg_conn"`
}
