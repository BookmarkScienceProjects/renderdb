package geometry

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
	"github.com/jmoiron/sqlx"
	"github.com/ungerik/go3d/vec3"
)

// Repository represents a spatial database with fast spatial lookups.
type Repository interface {
	// Add puts the object given in the database. Returns the ID of the inserted
	// object or an error.
	Add(o GeometryObject) (int64, error)
	// GetInsideVolume returns all objects inside the bounding box. Returns an error
	// if the database lookup fails (but not if the result set is empty).
	GetInsideVolume(bounds vec3.Box) ([]GeometryObject, error)
}

// NewRepository initializes a new repository using the given database.
func NewRepository(db *sqlx.DB) (Repository, error) {
	repo := new(defaultRepository)
	repo.database = newSQLGeometryDatabase(db)
	repo.tree = rtreego.NewTree(3, 25, 50)
	return repo, nil
}

type defaultRepository struct {
	database geometryDatabase
	tree     *rtreego.Rtree
}

func boxToRect(box vec3.Box) *rtreego.Rect {
	min := box.Min
	lengths := vec3.Sub(&box.Max, &min)

	p0 := rtreego.Point{float64(min[0]), float64(min[1]), float64(min[2])}
	l := rtreego.Point{float64(lengths[0]), float64(lengths[1]), float64(lengths[2])}
	rect, _ := rtreego.NewRect(p0, l)
	return rect
}

func rectToBox(rect *rtreego.Rect) vec3.Box {
	min := vec3.T{float32(rect.PointCoord(0)), float32(rect.PointCoord(1)), float32(rect.PointCoord(2))}
	lengths := vec3.T{float32(rect.LengthsCoord(0)), float32(rect.LengthsCoord(1)), float32(rect.LengthsCoord(2))}
	max := vec3.Add(&min, &lengths)
	return vec3.Box{min, max}
}

func (r *defaultRepository) Add(o GeometryObject) (int64, error) {
	id, err := r.database.add(o)
	if err == nil {
		r.tree.Insert(&rtreeEntry{id, boxToRect(o.Bounds())})
	}
	return id, err
}

func (r *defaultRepository) GetInsideVolume(bounds vec3.Box) ([]GeometryObject, error) {
	// Spacial lookup
	results := r.tree.SearchIntersect(boxToRect(bounds))
	ids := make([]int64, len(results))
	geometry := make(map[int64]*SimpleGeometryObject)
	for i, x := range results {
		entry := x.(*rtreeEntry)
		ids[i] = entry.id

		o := new(SimpleGeometryObject)
		o.bounds = rectToBox(entry.bounds)
		geometry[entry.id] = o
	}

	// Lookup exact geometry and metadata
	data, err := r.database.getMany(ids)
	if err != nil {
		return nil, err
	}

	// Merge spatial data and metadata/exact geometry
	for _, x := range data {
		o, found := geometry[x.id]
		if !found {
			return nil, fmt.Errorf("Database returned item with ID %d, but this was not in the query volume", x.id)
		}
		o.metadata = x.metadata
		o.geometryText = x.geometryText
	}

	// Extract
	asArray := make([]GeometryObject, 0, len(geometry))
	for _, v := range geometry {
		asArray = append(asArray, v)
	}
	return asArray, nil
}
