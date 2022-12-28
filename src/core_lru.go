package src

// Not support type-safe
type CoreLRU[K comparable, V any] struct {
	list  *list[K]
	items map[K]cacheItem[K, V]
	cap   int
}

func NewLRU[K comparable, V any](cap int) *CoreLRU[K, V] {
	if cap <= 0 {
		cap = 1
	}
	return &CoreLRU[K, V]{
		list:  newList[K](),
		items: make(map[K]cacheItem[K, V]),
		cap:   cap,
	}
}

// Return true if an empty is evited
func (cache *CoreLRU[K, V]) Add(key K, value V) (evicted bool) {
	item, ok := cache.items[key]

	if ok {
		item.value = value
		cache.list.moveToFront(item.node)
	}

	if cache.Len() >= cache.cap {
		e := cache.list.removeLast()
		delete(cache.items, e.value)
		evicted = true
	}

	ele := &node[K]{value: key}
	cache.list.push(ele)
	cache.items[key] = cacheItem[K, V]{node: ele, value: value}
	return true
}

func (cache *CoreLRU[K, V]) Len() int {
	return len(cache.items)
}

// Peek return the value without mark the key is latest
func (cache *CoreLRU[K, V]) Peek(key K) (V, bool) {
	item, ok := cache.items[key]
	return item.value, ok
}

func (cache *CoreLRU[K, V]) Get(key K) (V, bool) {
	item, ok := cache.items[key]

	if !ok {
		return item.value, ok
	}

	cache.list.moveToFront(item.node)
	return item.value, ok
}

func (cache *CoreLRU[K, V]) Remove(key K) bool {
	item, ok := cache.items[key]
	if ok {
		delete(cache.items, key)
		cache.list.remove(item.node)
	}
	return ok
}

func (cache *CoreLRU[K, V]) Purge() {
	cache.list = newList[K]()
	for k := range cache.items {
		delete(cache.items, k)
	}
}

type cacheItem[K any, V any] struct {
	node  *node[K]
	value V
}

type list[T any] struct {
	rootNode *node[T]
}

type node[T any] struct {
	next     *node[T]
	previous *node[T]
	value    T
}

func newList[T any]() *list[T] {
	root := node[T]{}
	root.next = &root
	root.previous = &root

	l := list[T]{rootNode: &root}
	return &l
}

// push to front of the list
func (l *list[T]) push(e *node[T]) {
	e.previous = l.rootNode
	e.next = l.rootNode.next
	l.rootNode.next = e
	e.next.previous = e
}

func (l *list[T]) moveToFront(e *node[T]) {
	e.next.previous = e.previous
	e.previous.next = e.next
	l.push(e)
}

func (l *list[T]) remove(e *node[T]) {
	e.previous.next = e.next
	e.next.previous = e.previous
	e.next = nil
	e.previous = nil
}

func (l *list[T]) removeLast() *node[T] {
	e := l.last()
	if e != nil {
		l.remove(e)
	}
	return e
}

func (l *list[T]) last() *node[T] {
	e := l.rootNode.previous

	if e == l.rootNode {
		return nil
	}
	return e
}
