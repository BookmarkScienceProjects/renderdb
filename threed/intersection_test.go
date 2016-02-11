package threed

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestRayTriangleIntersects_UnitTriangle_IntersectingRay_ReturnsTrue(t *testing.T) {
	// Arrange
	v1, v2, v3 := vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{0, 1, 0}
	o, d := vec3.T{0.2, 0.2, -1}, vec3.T{0, 0, 1}

	// Act
	intersects := RayTriangleIntersects(&v1, &v2, &v3, &o, &d)

	// Assert
	assert.True(t, intersects)
}

func TestRayTriangleIntersects_UnitTriangle_NonIntersectingRay_ReturnsFalse(t *testing.T) {
	// Arrange
	v1, v2, v3 := vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{0, 1, 0}
	o, d := vec3.T{0.2, 0.2, -1}, vec3.T{1, 1, 1}

	// Act
	intersects := RayTriangleIntersects(&v1, &v2, &v3, &o, &d)

	// Assert
	assert.False(t, intersects)
}

func TestRayTriangleIntersects_DegeneratedTriangle_DoesNotPanic(t *testing.T) {
	// Arrange
	v1, v2, v3 := vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{1, 0, 0}
	o, d := vec3.T{0.2, 0.2, -1}, vec3.T{0, 0, 1}

	// Act & Assert
	assert.NotPanics(t, func() { RayTriangleIntersects(&v1, &v2, &v3, &o, &d) })
	assert.NotPanics(t, func() { RayTriangleIntersects(&v2, &v1, &v3, &o, &d) })
	assert.NotPanics(t, func() { RayTriangleIntersects(&v3, &v1, &v2, &o, &d) })
}
