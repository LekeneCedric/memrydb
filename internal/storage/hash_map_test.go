package storage

import (
	"bytes"
	"testing"
)

const (
	default_shared_map_size = 8
)

func TestSharedMapCanStoreAndRetrieveValue(t *testing.T) {
	sm := NewSharedMap(default_shared_map_size)
	var value []byte = []byte("{\"name\": \"cedric\", \"surname\": \"michael\"}")

	sm.Set("user:1000", value)

	if got := sm.Get("user:1000"); !bytes.Equal(value, got) {
		t.Errorf("Incorrect value from SharedMap , wait for %s, but received %s", value, got)
	}
}

func TestSharedMapReturnNilIfKeyNotFounded(t *testing.T) {
	sm := NewSharedMap(default_shared_map_size)

	if got := sm.Get("user:1000"); got != nil {
		t.Error("An not founded key doesn't return an empty value")
	}
}
