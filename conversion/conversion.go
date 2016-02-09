// Package conversion contains utilities for doing basic type conversion.
package conversion

import (
	"github.com/dhconnelly/rtreego"
	"github.com/ungerik/go3d/float64/vec3"
)

// BoxToRect converts from vec3.Box to rtreego.Rect
func BoxToRect(box *vec3.Box) *rtreego.Rect {
	min := box.Min
	lengths := vec3.Sub(&box.Max, &min)

	p0 := rtreego.Point{min[0], min[1], min[2]}
	l := rtreego.Point{lengths[0], lengths[1], lengths[2]}
	rect, _ := rtreego.NewRect(p0, l)
	return rect
}

// RectToBox converts from rtreego.Rect to vec3.Box.
func RectToBox(rect *rtreego.Rect) *vec3.Box {
	min := vec3.T{rect.PointCoord(0), rect.PointCoord(1), rect.PointCoord(2)}
	lengths := vec3.T{rect.LengthsCoord(0), rect.LengthsCoord(1), rect.LengthsCoord(2)}
	max := vec3.Add(&min, &lengths)
	return &vec3.Box{min, max}
}

// SpatialSliceToBoundsSlice converts from []rtreego.Spatial to []*vec3.Box.
func SpatialSliceToBoundsSlice(objects []rtreego.Spatial) []*vec3.Box {
	// Spatial -> vec3.Box
	bounds := make([]*vec3.Box, len(objects))
	for i, s := range objects {
		bounds[i] = RectToBox(s.Bounds())
	}
	return bounds
}
