package config

type Config struct {
  NodeName          string      `json:"node_name"`
  RpcAddr           string      `json:"rpc_addr"`

  LogPath           string      `json:"log_path"`
  DBPath            string      `json:"db_path"`
  DataPath          string      `json:"data_path"`
}
