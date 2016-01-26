package options

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
	"github.com/larsmoa/renderdb/conversion"
	"github.com/ungerik/go3d/vec3"
)

// FilterGeometryOption is an interface to operations that filter/reshuffled
// elements based on geometry (bounds). This can e.g. be operations that
// cut of small objects far away or sorting objects by distance.
type FilterGeometryOption interface {
	// Apply filters or reshuffles elements by looking at the bounds.
	// Returns indices relative to the collection of bounds of the filtered/reshuffled
	// elements.
	Apply(bounds []*vec3.Box) []int
}

// VerifyAllAreOptions checks that all provided arguments are valid options.
func VerifyAllAreOptions(opts ...interface{}) error {
	for i, o := range opts {
		_, isOption := o.(FilterGeometryOption)
		if !isOption {
			return fmt.Errorf("Argument %d (%v) is not a valid option", i, o)
		}
	}
	return nil
}

// ApplyAllFilterGeometryOptions applies all FilterGeometryOption from the opts provided
// on the unfiltered objects-list given. The filters are applied in the order they are provided.
func ApplyAllFilterGeometryOptions(objects []rtreego.Spatial, opts ...interface{}) []rtreego.Spatial {
	for _, o := range opts {
		if option, ok := o.(FilterGeometryOption); ok {
			bounds := conversion.SpatialSliceToBoundsSlice(objects)
			keptIndices := option.Apply(bounds)

			// Apply filter/shuffle
			keep := make([]rtreego.Spatial, len(keptIndices))
			for i, j := range keptIndices {
				keep[i] = objects[j]
			}
			objects = keep
		}
	}

	return objects
}
