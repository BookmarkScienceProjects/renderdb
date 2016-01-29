package geometry

import (
	"fmt"
	"log"

	"github.com/dhconnelly/rtreego"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/conversion"
	"github.com/larsmoa/renderdb/geometry/options"
	"github.com/ungerik/go3d/vec3"
)

// Repository represents a spatial database with fast spatial lookups.
type Repository interface {
	// Add puts the object given in the database. Returns the ID of the inserted
	// object or an error.
	Add(o Object) (int64, error)
	// GetInsideVolume returns all objects inside the bounding box. Returns two channels,
	// one for geometry object and one for error. The operation is aborted on the first error.
	// Optionally, one or more Options may be provided to alter the behaviour of the
	// operation.
	GetInsideVolume(bounds vec3.Box, options ...interface{}) (<-chan Object, <-chan error)
}

// NewRepository initializes a new repository using the given database.
func NewRepository(db *sqlx.DB) (Repository, error) {
	repo := new(defaultRepository)
	repo.database = newSQLDatabase(db)
	repo.tree = rtreego.NewTree(3, 25, 50)
	if err := repo.loadFromDatabase(); err != nil {
		return nil, err
	}
	return repo, nil
}

type defaultRepository struct {
	database database
	tree     *rtreego.Rtree
}

func (r *defaultRepository) loadFromDatabase() error {
	dataCh, errCh := r.database.getAll()
	more := true
	log.Println("Initializing geometry database...")
	for more {
		var err error
		var d *data
		select {
		case d, more = <-dataCh:
			if more {
				treeEntry := new(rtreeEntry)
				treeEntry.id = d.id
				treeEntry.bounds = conversion.BoxToRect(&d.bounds)
				log.Printf("%d: %v\n", treeEntry.id, treeEntry.bounds)
				r.tree.Insert(treeEntry)
			}
		case err, more = <-errCh:
			if more {
				return err
			}
		}
	}
	log.Printf("Loaded %d geometry objects from database\n", r.tree.Size())
	return nil
}

func (r *defaultRepository) Add(o Object) (int64, error) {
	id, err := r.database.add(o)
	if err == nil {
		r.tree.Insert(&rtreeEntry{id, conversion.BoxToRect(o.Bounds())})
	}
	return id, err
}

func (r *defaultRepository) GetInsideVolume(bounds vec3.Box, opts ...interface{}) (<-chan Object, <-chan error) {
	geometryCh := make(chan Object, 200)
	errCh := make(chan error)

	go func() {
		defer close(geometryCh)

		// Verify arguments
		err := options.VerifyAllAreOptions(opts...)
		if err != nil {
			errCh <- err
			return
		}

		// Spacial lookup
		results := r.tree.SearchIntersect(conversion.BoxToRect(&bounds))

		// Apply geometry filters
		results = options.ApplyAllFilterGeometryOptions(results, opts...)

		// Create 'object placeholders'
		ids := make([]int64, len(results))
		geometry := make(map[int64]*SimpleObject)
		for i, x := range results {
			entry := x.(*rtreeEntry)
			ids[i] = entry.id

			o := new(SimpleObject)
			o.bounds = conversion.RectToBox(entry.bounds)
			geometry[entry.id] = o
		}

		// Lookup exact geometry and metadata
		dbDataCh, dbErrCh := r.database.getMany(ids)
		// Merge spatial data and metadata/exact geometry
		open := true
		for open {
			var data *data
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
