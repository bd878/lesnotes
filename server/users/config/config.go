package config

type Config struct {
	NodeName               string          `json:"node_name"`
	HttpAddr               string          `json:"http_addr"`
	SessionsServiceAddr    string          `json:"sessions_service_addr"`
	UsersServiceAddr       string          `json:"users_service_addr"`
	LogLevel               string          `json:"log_level"`
	NatsAddr               string          `json:"nats_addr"`
	NatsStream           string      `json:"nats_stream"`
	CookieDomain           string          `json:"cookie_domain"`
}
