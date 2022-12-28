package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoreLRU(t *testing.T) {
	lru := NewLRU[int, int](128)

	lru.Add(1, 1)
	value, ok := lru.Get(1)

	assert.Equal(t, true, ok)
	assert.Equal(t, 1, value)

	assert.Equal(t, 1, lru.Len())

	ok = lru.Remove(1)
	assert.Equal(t, true, ok)

	value, ok = lru.Get(1)
	assert.Equal(t, false, ok)

	assert.Equal(t, 0, lru.Len())
}

func TestCoreLRUMoreCap(t *testing.T) {
	lru := NewLRU[int, int](2)

	lru.Add(1, 1)
	lru.Add(2, 2)
	lru.Add(3, 3)

	value, ok := lru.Get(1)
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, value)
}
