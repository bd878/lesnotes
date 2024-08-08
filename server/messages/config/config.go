package config

type Config struct {
  NodeName          string `json:"node_name"`
  HttpAddr          string `json:"http_addr"`
  GrpcAddr          string `json:"grpc_addr"`
  LeaderAddr        string `json:"leader_addr"`
  UsersServiceAddr  string `json:"users_service_addr"`
  RaftAddr          string `json:"raft_addr"`
  SerfAddr          string `json:"serf_addr"`
  RaftServers       []string `json:"raft_servers"`
  SerfJoinAddrs     []string `json:"serf_join_addrs"`
  RaftLogLevel      string `json:"raft_log_level"`

  Debug             bool `json:"debug"`
  Bootstrap         bool `json:"bootstrap"`
  LogFile           string `json:"log_file"`
  DBPath            string `json:"db_path"`
  DataPath          string `json:"data_path"`
}
