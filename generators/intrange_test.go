package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntRange(t *testing.T) {
	assert.Empty(t, IntRange(0, 0))
	assert.EqualValues(t, []int{0, 1, 2, 3}, IntRange(0, 4))
	assert.EqualValues(t, []int{-5, -4, -3, -2}, IntRange(-5, 4))
	assert.EqualValues(t, []int{5, 6, 7}, IntRange(5, 3))
	assert.Panics(t, func() {
		IntRange(0, -10)
	})
}
