package repository

import (
	"log"

	"github.com/larsmoa/renderdb/conversion"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/repository/options"

	"github.com/dhconnelly/rtreego"
	"github.com/ungerik/go3d/float64/vec3"
)

// Repository represents a spatial database with fast spatial lookups.
type Repository interface {
	// Add puts the object given in the database. Returns the ID of the inserted
	// object or an error.
	Add(o db.Object) (int64, error)
	// GetInsideVolume returns all objects inside the bounding box. Returns two channels,
	// one for geometry object and one for error. The operation is aborted on the first error.
	// Optionally, one or more Options may be provided to alter the behaviour of the
	// operation.
	GetInsideVolume(bounds vec3.Box, options ...interface{}) (<-chan db.Object, <-chan error)
	// GetInsideVolumeIDs returns the same result as GetInsideVolume, but only returns
	// object IDs as a flat array rather than a channel of objects.
	GetInsideVolumeIDs(bounds vec3.Box, options ...interface{}) ([]int64, error)
	// GetWithIds returns objects with the given IDs. Returns two channels,
	// one for geometry object and one for error. The operation is aborted on the first error.
	GetWithIDs(ids []int64) (<-chan db.Object, <-chan error)
	// GetWithID returns object the the given ID, or an error if the operation fails.
	GetWithID(id int64) (db.Object, error)
}

/*
// NewRepository initializes a new repository using the given database.
func NewRepository(database *sqlx.DB) (Repository, error) {
	tx, err := database.Beginx()
	if err != nil {
		return nil, err
	}

	repo := new(defaultRepository)
	repo.database = db.NewObjectsDb(tx)
	repo.tree = rtreego.NewTree(3, 25, 50)
	if err := repo.loadFromDatabase(); err != nil {
		return nil, err
	}
	return repo, nil
}
*/

func NewRepository(database db.Objects) (Repository, error) {
	repo := new(defaultRepository)
	//repo.database = db.NewObjectsDb(tx)
	repo.database = database
	repo.tree = rtreego.NewTree(3, 25, 50)
	if err := repo.loadFromDatabase(); err != nil {
		return nil, err
	}
	return repo, nil
}

type defaultRepository struct {
	database db.Objects
	tree     *rtreego.Rtree
}

func (r *defaultRepository) loadFromDatabase() error {
	dataCh, errCh := r.database.GetAll()
	more := true
	log.Println("Initializing geometry database...")
	for more {
		var err error
		var d db.Object
		select {
		case d, more = <-dataCh:
			if more {
				treeEntry := new(rtreeEntry)
				treeEntry.id = d.ID()
				treeEntry.bounds = conversion.BoxToRect(d.Bounds())
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

func (r *defaultRepository) Add(o db.Object) (int64, error) {
	id, err := r.database.Add(o)
	if err == nil {
		r.tree.Insert(&rtreeEntry{id, conversion.BoxToRect(o.Bounds())})
	}
	return id, err
}

func (r *defaultRepository) GetInsideVolume(bounds vec3.Box, opts ...interface{}) (<-chan db.Object, <-chan error) {
	geometryCh := make(chan db.Object, 200)
	errCh := make(chan error)

	go func() {
		defer close(geometryCh)

		// Find IDs
		ids, err := r.GetInsideVolumeIDs(bounds, opts...)
		if err != nil {
			errCh <- err
			return
		}

		// Lookup exact geometry and metadata
		r.retrieveGeometryFromDatabase(ids, geometryCh, errCh)
	}()

	return geometryCh, errCh
}

func (r *defaultRepository) GetInsideVolumeIDs(bounds vec3.Box, opts ...interface{}) ([]int64, error) {
	// Verify arguments
	err := options.VerifyAllAreOptions(opts...)
	if err != nil {
		return nil, err
	}

	// Spacial lookup
	results := r.tree.SearchIntersect(conversion.BoxToRect(&bounds))

	// Apply geometry filters
	results = options.ApplyAllFilterGeometryOptions(results, opts...)

	// Extract IDs
	ids := make([]int64, len(results))
	for i, x := range results {
		entry := x.(*rtreeEntry)
		ids[i] = entry.id
	}

	return ids, nil
}

func (r *defaultRepository) GetWithIDs(ids []int64) (<-chan db.Object, <-chan error) {
	geometryCh := make(chan db.Object, 200)
	errCh := make(chan error)
	go func() {
		defer close(geometryCh)

		r.retrieveGeometryFromDatabase(ids, geometryCh, errCh)
	}()
	return geometryCh, errCh
}

func (r *defaultRepository) GetWithID(id int64) (db.Object, error) {
	geometryCh, errCh := r.GetWithIDs([]int64{id})
	select {
	case geom := <-geometryCh:
		return geom, nil
	case err := <-errCh:
		return nil, err
	}
}

func (r *defaultRepository) retrieveGeometryFromDatabase(ids []int64, geometryCh chan db.Object, errCh chan error) {
	if len(ids) == 0 {
		return
	}
	// Lookup exact geometry and metadata
	dbDataCh, dbErrCh := r.database.GetMany(ids)
	// Merge spatial data and metadata/exact geometry
	open := true
	for open {
		var object db.Object
		var err error
		select {
		case object, open = <-dbDataCh:
			if open {
				geometryCh <- object
			}

		case err, open = <-dbErrCh:
			if open {
				errCh <- err
			}
			return
		}
	}
}
