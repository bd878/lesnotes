package config

type Config struct {
	NodeName          string      `json:"node_name"`
	HttpAddr          string      `json:"http_addr"`
	LogPath           string      `json:"log_path"`
}
