package formats

import "github.com/ungerik/go3d/float64/vec3"

// geometryGroupAdapter implements the GeometryGroup interface
// given an objBuffer and a group.
type geometryGroupAdapter struct {
	buffer *objBuffer
	g      group
}

func createGeometryGroupAdapter(buffer *objBuffer, g group) *geometryGroupAdapter {
	groupBuffer := g.buildBuffers(buffer)
	return &geometryGroupAdapter{
		g:      g,
		buffer: groupBuffer,
	}
}

func (a *geometryGroupAdapter) Name() string {
	return a.g.name
}

func (a *geometryGroupAdapter) BoundingBox() vec3.Box {
	return a.buffer.BoundingBox()
}
