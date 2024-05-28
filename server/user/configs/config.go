package configs

type Config struct {
  HttpPort int `json:"httpport"`
  GrpcPort int `json:"grpcport"`
  Debug bool `json:"debug"`
  LogFile string `json:"logFile"`
  DBPath string `json:"dbPath"`
}
