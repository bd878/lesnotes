package main

type config struct {
  Port int `json:"port"`
  Debug bool `json:"debug"`
  LogFile string `json:"logFile"`
}
