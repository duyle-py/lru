package src

import (
	ll "container/list"
	"math"
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
	cache map[K]*page[K, V]
}

type page[K comparable, V comparable] struct {
	Key     K
	Value   V
	TAdress *ll.List
	Ele     *ll.Element
}

func New[K comparable, V comparable](size int) *cacheARC[K, V] {
	return &cacheARC[K, V]{
		p:     0,
		size:  size,
		t1:    ll.New(),
		b1:    ll.New(),
		t2:    ll.New(),
		b2:    ll.New(),
		cache: make(map[K]*page[K, V], size),
	}
}

// Get Value from cache by key
func (c *cacheARC[K, V]) Get(key K) *V {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	p, ok := c.cache[key]

	if !ok {
		return nil
	}

	return &p.Value

}

// Add inserts a new key-value pair into the cache.
func (c *cacheARC[K, V]) Add(key K, val V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	p, ok := c.cache[key]
	if !ok {
		p = &page[K, V]{
			Key:   key,
			Value: val,
		}
		c.req(p)
		c.cache[key] = p
	} else {
		p.Value = val
		c.req(p)
	}
	return ok
}

// Element request
// Case 1:
//
//	page in t1 || t2
//
// Case 2:
//
//	page in b1
//
// Case 3:
//
//	page in b2
//
// Case 4:
//
//	page not in t1, t2, b1, b2
func (c *cacheARC[K, V]) req(p *page[K, V]) {
	// Case I
	if p.TAdress == c.t1 || p.TAdress == c.t2 {
		p.setMRU(c.t2)
	} else if p.TAdress == c.b1 {
		// Case II not in t1, t2
		// Adaptation
		c.p = int(math.Min(float64(c.size), float64(float64(c.p)+math.Max(float64(c.b2.Len()/c.b1.Len()), float64(1)))))
		c.replace(p)
		p.setMRU(c.t2)
	} else if p.TAdress == c.b2 {
		// Case III not in t1, t2
		// Adaptation
		c.p = int(math.Max(float64(0), float64(float64(c.p)-math.Max(float64(c.b1.Len()/c.b2.Len()), float64(1)))))
		c.replace(p)
		p.setMRU(c.t2)
	} else if p.TAdress == nil {
		// Case IV
		// Case 1
		if c.t1.Len()+c.b1.Len() == c.size {
			if c.t1.Len() < c.size {
				delLRU(c.b1)
				c.replace(p)
			} else {
				lru := delLRU(c.t1)
				delete(c.cache, lru.Value.(*page[K, V]).Key)
			}

		} else if c.t1.Len()+c.b1.Len() < c.size && c.t1.Len()+c.b1.Len()+c.t2.Len()+c.b2.Len() >= c.size {
			if c.t1.Len()+c.b1.Len()+c.t2.Len()+c.b2.Len() >= 2*c.size {
				_ = delLRU(c.b2)
			}
			c.replace(p)
		}
		p.setMRU(c.t1)
	}
}

// Replace function
func (c *cacheARC[K, V]) replace(p *page[K, V]) {
	if c.t1.Len() >= 1 && ((p.TAdress == c.b2 && c.t1.Len() == c.p) || (c.t1.Len() > c.p)) {
		// Move page from t1 to b1 and remove it in cache
		lru := delLRU(c.t1).Value.(*page[K, V])
		lru.setMRU(c.b1)
		delete(c.cache, lru.Key)
	} else {
		// Move page from t2 to b2 and remove it in cache
		lru := delLRU(c.t2).Value.(*page[K, V])
		lru.setMRU(c.b2)
		delete(c.cache, lru.Key)
	}
}

// Del LRU
func delLRU(ll *ll.List) *ll.Element {
	lru := ll.Back()
	ll.Remove(lru)
	return lru
}

func (p *page[K, V]) setLRU(l *ll.List) {
	p.detach()
	p.TAdress = l
	p.Ele = p.TAdress.PushBack(p)
}

func (p *page[K, V]) setMRU(l *ll.List) {
	p.detach()
	p.TAdress = l
	p.Ele = p.TAdress.PushFront(p)
}

func (p *page[K, V]) detach() {
	if p.TAdress != nil {
		p.TAdress.Remove(p.Ele)
	}
}
