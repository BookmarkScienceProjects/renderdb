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
	Add(o Object) (int64, error)
	// GetInsideVolume returns all objects inside the bounding box. Returns two channels,
	// one for geometry object and one for error. The operation is aborted on the first error.
	GetInsideVolume(bounds vec3.Box) (<-chan Object, <-chan error)
}

// NewRepository initializes a new repository using the given database.
func NewRepository(db *sqlx.DB) (Repository, error) {
	repo := new(defaultRepository)
	repo.database = newSQLDatabase(db)
	repo.tree = rtreego.NewTree(3, 25, 50)
	return repo, nil
}

type defaultRepository struct {
	database database
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

func (r *defaultRepository) Add(o Object) (int64, error) {
	id, err := r.database.add(o)
	if err == nil {
		r.tree.Insert(&rtreeEntry{id, boxToRect(o.Bounds())})
	}
	return id, err
}

func (r *defaultRepository) GetInsideVolume(bounds vec3.Box) (<-chan Object, <-chan error) {
	geometryCh := make(chan Object, 200)
	errCh := make(chan error)

	go func() {
		defer close(geometryCh)

		// Spacial lookup
		results := r.tree.SearchIntersect(boxToRect(bounds))
		ids := make([]int64, len(results))
		geometry := make(map[int64]*SimpleObject)
		for i, x := range results {
			entry := x.(*rtreeEntry)
			ids[i] = entry.id

			o := new(SimpleObject)
			o.bounds = rectToBox(entry.bounds)
			geometry[entry.id] = o
		}

		// Lookup exact geometry and metadata
		dbDataCh, dbErrCh := r.database.getMany(ids)
		// Merge spatial data and metadata/exact geometry
		open := true
		for open {
			var data *data
			var err error

			select {
			case data, open = <-dbDataCh:
				if open {
					o, found := geometry[data.id]
					if !found {
						errCh <- fmt.Errorf("Database returned item with ID %d, but this was not in the query volume", data.id)
						return
					}
					o.metadata = data.metadata
					o.geometryText = data.geometryText
					geometryCh <- o
				}

			case err, open = <-dbErrCh:
				if open {
					errCh <- err
				}
				return
			}
		}
	}()

	return geometryCh, errCh
}
