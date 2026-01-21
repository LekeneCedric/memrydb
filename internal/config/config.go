package config

import (
	"encoding/json"
	"io"
)

type Config struct {
	NumberOfShard        uint8  `json:"number_of_shard"`
	MaxMemoryUsage       uint16 `json:"max_memory_usage"`
	MaxNumberOfConnexion uint64 `json:"max_number_of_connexion"`
}

func NewConfig(r io.Reader) (*Config, error) {
	config := &Config{}
	err := json.NewDecoder(r).Decode(config)
	return config, err
}
