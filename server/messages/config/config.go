package config

type Config struct {
  PublicIp        string `json:"public_ip"`
  PrivateIp       string `json:"private_ip"`

  HttpPort        int `json:"http_port"`
  GrpcPort        int `json:"grpc_port"`
  LeaderPort      int `json:"leader_port"`
  UserPort        int `json:"user_port"`
  RaftStreamPort  int `json:"raft_stream_port"`

  Debug           bool `json:"debug"`
  LogFile         string `json:"log_file"`
  DBPath          string `json:"db_path"`
  DataPath        string `json:"data_path"`
  Bootstrap       bool `json:"bootstrap"`
}
