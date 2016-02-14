package formats

import (
	"errors"
	"fmt"

	"github.com/larsmoa/renderdb/threed"

	"github.com/ungerik/go3d/float64/vec3"
)

type lineError struct {
	lineNumber int
	line       string
	err        error
}

func (e lineError) Error() string {
	return fmt.Sprintf("Line #%d: %v ('%s')", e.lineNumber, e.line, e.err)
}

// faceCorner represents a 'corner' (or vertex) in a face
type faceCorner struct {
	vertexIndex int
	normalIndex int
}

// face represents a surface represented by a set of corner
type face struct {
	corners  []faceCorner
	material string
}

type objBuffer struct {
	activeMaterial string

	// All the below maps directly to OBJ-keywords
	mtllib string
	v      []vec3.T
	vn     []vec3.T
	f      []face
	g      []group
}

func (b *objBuffer) BoundingBox() vec3.Box {
	box := vec3.Box{vec3.MaxVal, vec3.MinVal}
	for _, v := range b.v {
		box.Join(&vec3.Box{v, v})
	}
	return box
}

func (b *objBuffer) Intersects(origin, direction *vec3.T) bool {
	for _, f := range b.f {
		if len(f.corners) != 3 {
			panic(errors.New("Intersects only works on triangles"))
		}

		v1 := &b.v[f.corners[0].vertexIndex]
		v2 := &b.v[f.corners[1].vertexIndex]
		v3 := &b.v[f.corners[2].vertexIndex]
		if threed.RayTriangleIntersects(v1, v2, v3, origin, direction) {
			return true
		}
	}
	return false
}

// ReadOptions represents options used by WavefrontObjReader.Read.
type ReadOptions struct {
	// DiscardDegeneratedFaces instructs the reader to discard faces
	// where a vertex index appears twice, e.g. "1 1 2".
	DiscardDegeneratedFaces bool
}
