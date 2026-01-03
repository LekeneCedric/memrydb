package storage

import (
	"sync"

	"github.com/LekeneCedric/memrydb/internal/utils"
)

type Shard struct {
	mux  sync.RWMutex
	data map[string][]byte
}

type SharedMap struct {
	shards []Shard
}

func NewSharedMap(size uint8) *SharedMap {
	if size == 0 {
		size = 32
	}
	sm := &SharedMap{
		shards: make([]Shard, size),
	}
	for i := range sm.shards {
		sm.shards[i].data = make(map[string][]byte)
	}
	return sm
}

func (sm *SharedMap) Get(key string) []byte {
	shard := sm.getShard(key)
	shard.mux.RLock()
	defer shard.mux.RUnlock()
	return shard.data[key]
}

func (sm *SharedMap) Set(key string, value []byte) {
	shard := sm.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	shard.data[key] = value
}

func (sm *SharedMap) Remove(key string) {
	shard := sm.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	delete(shard.data, key)
}

func (sm *SharedMap) getShard(key string) *Shard {
	hashedKey := utils.Hash32(key)
	shardIndex := int(hashedKey % uint32(len(sm.shards)))
	return &sm.shards[shardIndex]
}
