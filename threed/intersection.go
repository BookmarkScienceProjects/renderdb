package threed

import (
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

const (
	epsilon float64 = 1e-8
)

func toPtr(v vec3.T) *vec3.T {
	return &v
}

// RayTriangleIntersects performs an intersection between a ray and a triangle.
// Algorithm from http://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-rendering-a-triangle/moller-trumbore-ray-triangle-intersection
func RayTriangleIntersects(v0, v1, v2 *vec3.T, orig, dir *vec3.T) bool {
	v0v1 := v1.Sub(v0)
	v0v2 := v2.Sub(v0)
	pvec := toPtr(vec3.Cross(dir, v0v2))
	det := vec3.Dot(v0v1, pvec)

	// Ray and triangle are parallel if det is close to 0
	if math.Abs(det) < epsilon {
		return false
	}

	invDet := 1 / det
	tvec := orig.Sub(v0)

	u := vec3.Dot(tvec, pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	qvec := toPtr(vec3.Cross(tvec, v0v1))
	v := vec3.Dot(dir, qvec) * invDet
	if v < 0 || u+v > 1 {
		return false
	}

	return true
}
