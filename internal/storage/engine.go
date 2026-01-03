package storage

type Engine interface {
	Get(key string) []byte
	Set(key string, value []byte)
	Remove(key string)
}
