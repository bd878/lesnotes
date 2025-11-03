package config

type Config struct {
	NodeName             string      `json:"node_name"`
	RpcAddr              string      `json:"rpc_addr"`
	HttpAddr             string      `json:"http_addr"`
	UsersServiceAddr     string      `json:"users_service_addr"`
	SessionsServiceAddr  string      `json:"sessions_service_addr"`

	PGConn               string      `json:"pg_conn"`
}
