package config

type Config struct {
	HttpAddr          string          `json:"http_addr"`
	RpcAddr           string          `json:"rpc_addr"`
	MessagesServiceAddr string        `json:"messages_service_addr"`
	NodeName          string          `json:"node_name"`
	LogPath           string          `json:"log_path"`
	DataPath          string          `json:"data_path"`
	DBPath            string          `json:"db_path"`
	CookieDomain      string          `json:"cookie_domain"`
}
