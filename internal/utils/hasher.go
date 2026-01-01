package utils

const (
	offset32 = 2166136261
	prime32  = 16777619
)

func Hash32(key string) uint32 {
	var hash uint32 = offset32
	var data []byte = []byte(key)

	for i := range data {
		hash ^= uint32(data[i])
		hash *= prime32
	}

	return hash
}
