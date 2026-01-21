package config

import (
	"bytes"
	"testing"
)

const (
	NUMBER_OF_SHARD         uint8  = 32
	MAX_MEMORY_USAGE        uint16 = 4026
	MAX_NUMBER_OF_CONNEXION uint64 = 1024
)

func Test_LoadConfigs(t *testing.T) {
	rawConfigFile := []byte(`{"number_of_shard":32,"max_memory_usage":4026,"max_number_of_connexion":1024}`)
	configs, _ := NewConfig(bytes.NewReader(rawConfigFile))

	type test struct {
		name     string
		expected any
		actual   any
	}

	tests := []test{
		{name: "number of shard", expected: NUMBER_OF_SHARD, actual: configs.NumberOfShard},
		{name: "max memory usage", expected: MAX_MEMORY_USAGE, actual: configs.MaxMemoryUsage},
		{name: "max number of connexion", expected: MAX_NUMBER_OF_CONNEXION, actual: configs.MaxNumberOfConnexion},
	}

	for _, tt := range tests {
		if tt.actual != tt.expected {
			t.Fatalf("config: %s : %d expected , %d received", tt.name, tt.expected, tt.actual)
		}
	}
}

func Test_RejectInvalidConfigsFile(t *testing.T) {
	rawConfigFile := []byte(`{"number_of_shard":32,"max_memory_usage":4026,"max_number_of_connexion":10`)
	_, err := NewConfig(bytes.NewReader(rawConfigFile))
	if err == nil {
		t.Fatalf("Should reject invalid configs file")
	}
}
