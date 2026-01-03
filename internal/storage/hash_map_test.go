package storage

import (
	"bytes"
	"testing"
)

const (
	default_shared_map_size = 8
	key                     = "user:1000"
)

func TestSharedMapCanStoreAndRetrieveValue(t *testing.T) {
	sm := NewSharedMap(default_shared_map_size)
	var value []byte = []byte("{\"name\": \"cedric\", \"surname\": \"michael\"}")

	sm.Set(key, value)

	if got := sm.Get(key); !bytes.Equal(value, got) {
		t.Errorf("Incorrect value from SharedMap , wait for %s, but received %s", value, got)
	}
}

func TestSharedMapReturnNilIfKeyNotFounded(t *testing.T) {
	sm := NewSharedMap(default_shared_map_size)

	if got := sm.Get(key); got != nil {
		t.Error("An not founded key doesn't return an empty value")
	}
}

func TestShareMapCanRemoveElement(t *testing.T) {
	//Given
	sm := NewSharedMap(default_shared_map_size)
	sm.Set(key, []byte("{\"name\": \"Lekene Cedric\"}"))
	//When
	sm.Remove(key)
	//Then
	if got := sm.Get(key); got != nil {
		t.Errorf("An empty output was expected , instead we get : %s", got)
	}
}
