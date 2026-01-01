package utils

import (
	"hash/fnv"
	"testing"
)

const (
	key = "user:100:session:active"
)

func BenchmarkCustomHash(b *testing.B) {
	b.ResetTimer()

	for b.Loop() {
		_ = Hash32(key)
	}
}

func BenchmarkStandardFNV(b *testing.B) {
	b.ResetTimer()

	for b.Loop() {
		h := fnv.New32()
		h.Write([]byte(key))
		_ = h.Sum32()
	}
}
