package src

import "sync"

// Safe for concurent usage
type Cache[K comparable, V any] struct {
	cache CoreLRU[K, V]
	mu    sync.Mutex
}
