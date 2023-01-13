package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestARC(t *testing.T) {
	arc := New[int, int](128)

	arc.Add(1, 1)
	value := arc.Get(1)

	assert.Equal(t, 1, *value)

	arc.Add(1, 3)
	value = arc.Get(1)

	assert.Equal(t, 3, *value)

	arc.Add(2, 4)
	value1 := arc.Get(2)

	assert.Equal(t, 3, *value)
	assert.Equal(t, 4, *value1)
}
