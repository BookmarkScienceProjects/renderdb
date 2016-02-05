package formats

import "github.com/ungerik/go3d/float64/vec3"

type objBuffer struct {
	facesets []faceset
	// All the below maps directly to OBJ-keywords
	mtllib string
	v      []vec3.T
	vn     []vec3.T
	f      []face
	g      []group
}
