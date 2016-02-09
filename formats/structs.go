package formats

import (
	"fmt"

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

// material represents the name of a material used by face.
type material string

const (
	// Special material representing an undefined material
	undefinedMaterial material = ""
)

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
