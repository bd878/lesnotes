package configs

type Config struct {
  FcgiPort int `json:"fcgiport"`
  GrpcPort int `json:"grpcport"`
  Debug bool `json:"debug"`
  LogFile string `json:"logFile"`
  DBPath string `json:"dbPath"`
}
