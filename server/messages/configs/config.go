package configs

type Config struct {
  Port int `json:"port"`
  UserAddr string `json:"useraddr`
  Debug bool `json:"debug"`
  LogFile string `json:"logFile"`
  DBPath string `json:"dbPath"`
  DataPath string `json:"dataPath"`
}
