package utils

import "testing"

const (
	input = "user:100"
)

func TestHasher_NewHash32_IsDeterministic(t *testing.T) {
	input1 := input
	input2 := input

	hash1 := Hash32(input1)
	hash2 := Hash32(input2)

	if hash1 != hash2 {
		t.Errorf("The hash function is not Deterministic, the same input should always return the same output : %d != %d", hash1, hash2)
	}
}
