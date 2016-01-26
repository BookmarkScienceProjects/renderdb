package options

import (
	"sort"

	"github.com/larsmoa/renderdb/generators"
	"github.com/ungerik/go3d/vec3"
)

type byDistance struct {
	pivot   vec3.T
	bounds  []*vec3.Box
	indices []int
}

func sqDistToClosestPointOnBox(p *vec3.T, box *vec3.Box) float32 {
	closestComponent := func(i int) float32 {
		if p[i] < box.Min[i] {
			return box.Min[i]
		} else if p[i] > box.Max[i] {
			return box.Max[i]
		} else {
			return p[i]
		}
	}
	// Find closest point on bounding box
	c := vec3.T{}
	for i := 0; i < 3; i++ {
		c[i] = closestComponent(i)
	}
	// Take distance between the two points
	return vec3.SquareDistance(p, &c)
}

func (d byDistance) Len() int {
	return len(d.indices)
}
func (d byDistance) Swap(i, j int) {
	d.indices[i], d.indices[j] = d.indices[j], d.indices[i]
}
func (d byDistance) Less(i, j int) bool {
	b1 := d.bounds[d.indices[i]]
	b2 := d.bounds[d.indices[j]]
	d1 := sqDistToClosestPointOnBox(&d.pivot, b1)
	d2 := sqDistToClosestPointOnBox(&d.pivot, b2)
	return d1 < d2
}

// SortByDistance sorts objects returned from e.g. Repository.GetInsideVolume()
// by distance (nearest first) to some 'pivot location'.
type SortByDistance struct {
	Pivot vec3.T
}

// Apply returns a list of indices that can be used to look up in
// bounds to create an ordered list.
func (o *SortByDistance) Apply(bounds []*vec3.Box) []int {
	byDistance := byDistance{o.Pivot, bounds, generators.IntRange(0, len(bounds))}
	sort.Sort(byDistance)
	return byDistance.indices
}
