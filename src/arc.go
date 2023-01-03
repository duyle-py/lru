package src

import (
	ll "container/list"
	"sync"
)

type cacheARC[K comparable, V comparable] struct {
	// Target size for t1
	p    int
	size int
	// Recent cache link list
	t1 *ll.List
	// Evicted from t1
	b1 *ll.List
	// Frequence Cache link list
	t2 *ll.List
	// Evicted from t2
	b2    *ll.List
	mutex sync.RWMutex
	len   int
	cache map[K]*ele[K, V]
}

type ele[K comparable, V comparable] struct {
	Key   K
	Value V
}

func New[K comparable, V comparable](size int) *cacheARC[K, V] {
	return &cacheARC[K, V]{
		p:     0,
		size:  size,
		t1:    ll.New(),
		b1:    ll.New(),
		t2:    ll.New(),
		b2:    ll.New(),
		len:   0,
		cache: make(map[K]*ele[K, V], size),
	}
}

// Get Value from cache by key
func (cache *cacheARC[K, V]) Get(key K) (V, bool)

//TODO

// Add to cache
func (cache *cacheARC[K, V]) Add(key K, val V) (V, bool)

//TODO

// Element request
// Case 1:
//
//	ele in t1 || t2
//		=> remove in t1 || t2
//		=> move to front t2
//
// Case 2:
//
//	ele in b1
//		=> TODO
//
// Case 3:
//
//	ele in b2
//		=> TODO
//
// Case 4:
//
//	ele not in t1, t2, b1, b2
//		=> TODO
func (cache *cacheARC[K, V]) req(val *ele[K, V])

// Replace function
func (cache *cacheARC[K, V]) replace(val *ele[K, V])
