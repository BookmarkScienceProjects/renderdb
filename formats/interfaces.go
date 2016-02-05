package formats

import (
	"io"

	"github.com/ungerik/go3d/float64/vec3"
)

type GeometryGroup interface {
	Name() string
	BoundingBox() vec3.Box
	Write(w io.Writer) error
}
