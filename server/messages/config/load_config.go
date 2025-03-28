package config

import (
	"os"
	"encoding/json"
)

func Load(configPath string) Config {
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	return cfg
}