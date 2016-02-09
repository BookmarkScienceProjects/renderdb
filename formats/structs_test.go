package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestObjBuffer_BoundingBox_NoVertices_ReturnsEmptyBox(t *testing.T) {
	// Arrange
	buffer := objBuffer{}

	// Act
	box := buffer.BoundingBox()

	// Assert
	assert.Equal(t, vec3.Box{vec3.MaxVal, vec3.MinVal}, box)
}

func TestObjBuffer_BoundingBox_WithVertices_ReturnsCorrectBoundingBox(t *testing.T) {
	// Arrange
	buffer := objBuffer{}
	buffer.v = []vec3.T{
		vec3.T{1, 0, 0}, vec3.T{0, 0, 0}, vec3.T{2, 3, 1},
		vec3.T{1, 1, 1}, vec3.T{1.5, 2.5, 4}, vec3.T{1, 1, 0},
	}

	// Act
	box := buffer.BoundingBox()

	// Assert
	assert.Equal(t, vec3.Box{Min: vec3.T{0, 0, 0}, Max: vec3.T{2, 3, 4}}, box)
}

func TestObjBuffer_BoundingBox_ZeroNotIncludedInBounds_ReturnsCorrectBoundingBox(t *testing.T) {
	// Arrange
	buffer := objBuffer{}
	buffer.v = []vec3.T{
		vec3.T{1, 1, 1}, vec3.T{2, 1, 3}, vec3.T{1, 4, 5},
	}

	// Act
	box := buffer.BoundingBox()

	// Assert
	assert.Equal(t, vec3.Box{Min: vec3.T{1, 1, 1}, Max: vec3.T{2, 4, 5}}, box)

}
