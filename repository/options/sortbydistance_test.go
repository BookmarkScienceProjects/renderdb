package options

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func createBox(center vec3.T, size float64) *vec3.Box {
	diff := vec3.T{size / 2.0, size / 2.0, size / 2.0}
	min := vec3.Sub(&center, &diff)
	max := vec3.Add(&center, &diff)
	return &vec3.Box{min, max}
}

func TestSortByDistance_TwoBoxesSameCenterDifferentSize_ReturnsCorrectOrder(t *testing.T) {
	// Arrange
	opt := SortByDistance{vec3.T{0, 0, 0}}
	bounds := []*vec3.Box{
		createBox(vec3.T{1, 1, 1}, 0.1),
		createBox(vec3.T{1, 1, 1}, 0.2),
	}

	// Act
	indices := opt.Apply(bounds)

	// Assert
	assert.EqualValues(t, []int{1, 0}, indices)
}

func TestSortByDistance_TwoBoxesDifferentCenterSameSize_ReturnsCorrectOrder(t *testing.T) {
	// Arrange
	opt := SortByDistance{vec3.T{0, 0, 0}}
	bounds := []*vec3.Box{
		createBox(vec3.T{1, 1, 1}, 1),
		createBox(vec3.T{1, 1, 2}, 1),
	}

	// Act
	indices := opt.Apply(bounds)

	// Assert
	assert.EqualValues(t, []int{0, 1}, indices)
}

func TestSortByDistance_BigBoxEncapsulatingPivot_IsReturnedFirst(t *testing.T) {
	// Arrange
	opt := SortByDistance{vec3.T{0, 0, 0}}
	bounds := []*vec3.Box{
		createBox(vec3.T{1, 1, 1}, 0.1),  // Doesn't contain pivot, but edge is close
		createBox(vec3.T{0, 0, 0}, 1000), // Contains pivot, but 'edge' is very far away
		createBox(vec3.T{1, 1, 1}, 0.9),  // Doesn't contain pivot, but edge is very close
	}

	// Act
	indices := opt.Apply(bounds)

	// Assert
	assert.EqualValues(t, []int{1, 2, 0}, indices)
}

func TestSortByDistance_1000RandomBoxesAroundPivot_DoesNotPanic(t *testing.T) {
	// Arrange
	opt := SortByDistance{vec3.T{0.5, 0.5, 0.5}}
	bounds := make([]*vec3.Box, 1000)
	for i := 0; i < len(bounds); i++ {
		bounds[i] = createBox(vec3.T{rand.Float64(), rand.Float64(), rand.Float64()}, rand.Float64())
	}

	// Act & Assert - test that all combinations of proximities 'work' (or doesn't crash at least)
	assert.NotPanics(t, func() { _ = opt.Apply(bounds) })
}
