package config

type Config struct {
  HttpAddr          string `json:"http_addr"`
  GrpcAddr          string `json:"grpc_addr"`
  LeaderAddr        string `json:"leader_addr"`
  UsersServiceAddr  string `json:"users_service_addr"`
  RaftAddr          string `json:"raft_addr"`
  JoinAddrs         []string `json:"join_addrs"`

  Debug             bool `json:"debug"`
  Bootstrap         bool `json:"bootstrap"`
  LogFile           string `json:"log_file"`
  DBPath            string `json:"db_path"`
  DataPath          string `json:"data_path"`
}
