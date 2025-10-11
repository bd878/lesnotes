package config

type Config struct {
	NodeName             string      `json:"node_name"`
	HttpAddr             string      `json:"http_addr"`
	UsersServiceAddr     string      `json:"users_service_addr"`
	SessionsServiceAddr  string      `json:"sessions_service_addr"`
	LogLevel             string      `json:"log_level"`

	NatsAddr             string      `json:"nats_addr"`

	PGConn               string      `json:"pg_conn"`
	MessagesTableName    string      `json:"messages_table_name"`
	FilesTableName       string      `json:"files_table_name"`
}
