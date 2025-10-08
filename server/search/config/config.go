package config

type Config struct {
	RpcAddr              string      `json:"rpc_addr"`
	HttpAddr             string      `json:"http_addr"`
	FilesServiceAddr     string      `json:"files_service_addr"`
	MessagesServiceAddr  string      `json:"messages_service_addr"`
	RpcAddr              string      `json:"rpc_addr"`
	LogLevel             string      `json:"log_level"`

	NatsAddr             string      `json:"nats_addr"`
	NatsStream           string      `json:"nats_stream"`

	PGConn               string      `json:"pg_conn"`
	TableName            string      `json:"table_name"`
}
