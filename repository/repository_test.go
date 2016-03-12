package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dhconnelly/rtreego"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/repository/options"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

// createGetManyResult is a helper function to emulate how getMany() returns
// values/error. Usage:
//   .Returns(createGetManyResult(data1, data2, err)) // Returns two values, then fails
func createGetManyResult(values ...interface{}) (<-chan db.Object, <-chan error) {
	dataCh := make(chan db.Object)
	errCh := make(chan error)
	go func() {
		defer close(dataCh)
		for _, v := range values {
			if data, ok := v.(db.Object); ok {
				dataCh <- data
			} else if err, ok := v.(error); ok {
				errCh <- err
				return
			} else {
				panic(fmt.Sprintf("Values must be error or db.Object, got %T", v))
			}
		}
	}()
	return dataCh, errCh
}

func flattenChannels(objCh <-chan db.Object, errCh <-chan error) ([]db.Object, error) {
	objects := []db.Object{}
	doBreak := false
	for !doBreak {
		select {
		case obj, more := <-objCh:
			if more {
				objects = append(objects, obj)
			} else {
				doBreak = true
			}
		case err := <-errCh:
			return objects, err
		}
	}
	return objects, nil
}

func TestRepository_Add_ValidGeometry_AddsToTreeAndDatabase(t *testing.T) {
	// Arrange
	obj := db.NewSimpleObject(vec3.Box{}, nil, nil)

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj).Return(int64(1), nil)

	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	id, err := repo.Add(obj)

	// Assert
	assert.Equal(t, int64(1), id)
	assert.Nil(t, err)
	assert.Equal(t, 1, rtree.Size())
	mockDb.AssertExpectations(t)
}

func TestRepository_Add_DatabaseReturnsError_DoesNotAddToTree(t *testing.T) {
	// Arrange
	obj := db.NewSimpleObject(vec3.Box{}, nil, nil)

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj).Return(int64(0), errors.New("error"))

	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	_, err := repo.Add(obj)

	// Assert
	assert.Equal(t, 0, rtree.Size())
	assert.NotNil(t, err)
}

func TestRepository_GetInsideVolume_NothingInsideVolume_ReturnsEmpty(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{1, 1, 1}, vec3.T{2, 2, 2}}
	obj := db.NewSimpleObject(objBounds, []byte{}, nil)

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj).Return(int64(1), nil)

	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	bounds := vec3.Box{vec3.T{5, 5, 5}, vec3.T{6, 6, 6}}
	objects, err := flattenChannels(repo.GetInsideVolume(bounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Empty(t, objects)
}

func TestRepository_GetInsideVolume_OneInsideVolume_ReturnsObject(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj := db.NewSimpleObject(objBounds, []byte{}, nil)

	data := new(db.MockObject)
	data.On("ID").Return(int64(1))

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj).Return(int64(1), nil)
	mockDb.On("GetMany", []int64{1}).Return(createGetManyResult(data))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := flattenChannels(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(objects))
}

func TestRepository_GetInsideVolume_DatabaseReturnsError_ReturnsError(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj := db.NewSimpleObject(objBounds, []byte{}, nil)

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj).Return(int64(1), nil)
	mockDb.On("GetMany", []int64{1}).Return(createGetManyResult(errors.New("error")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := flattenChannels(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, 0, len(objects))
}

func TestRepository_GetInsideVolume_DatabaseReturnsOneThenError_ReturnsError(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj1 := db.NewSimpleObject(objBounds, []byte("1"), nil)
	obj2 := db.NewSimpleObject(objBounds, []byte("2"), nil)

	data := new(db.MockObject)
	data.On("ID").Return(int64(1))

	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj1).Return(int64(1), nil)
	mockDb.On("Add", obj2).Return(int64(2), nil)
	mockDb.On("GetMany", []int64{1, 2}).Return(createGetManyResult(data, errors.New("error")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj1)
	repo.Add(obj2)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	_, err := flattenChannels(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.Error(t, err)
}

func TestRepository_GetInsideVolume_WithFilterGeometryOptions_ReturnsFiltered(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj1 := db.NewSimpleObject(objBounds, []byte("1"), nil)
	obj2 := db.NewSimpleObject(objBounds, []byte("2"), nil)

	data := new(db.MockObject)
	data.On("ID").Return(int64(1))
	mockDb := new(db.MockObjects)
	mockDb.On("Add", obj1).Return(int64(1), nil)
	mockDb.On("Add", obj2).Return(int64(2), nil)
	mockDb.On("GetMany", []int64{1}).Return(createGetManyResult(data))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj1)
	repo.Add(obj2)

	mockOptions := new(options.MockFilterGeometryOption)
	mockOptions.On("Apply", []*vec3.Box{&objBounds, &objBounds}).Return([]int{0})

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	result, err := flattenChannels(repo.GetInsideVolume(searchBounds, mockOptions))

	// Assert
	mockDb.AssertExpectations(t)
	mockOptions.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
}

func TestRepository_LoadFromDatabase_LoadsObjectsFromDatabase(t *testing.T) {
	// Arrange
	obj1 := new(db.MockObject)
	obj1.On("ID").Return(int64(1))
	obj1.On("Bounds").Return(&vec3.Box{})
	obj2 := new(db.MockObject)
	obj2.On("ID").Return(int64(2))
	obj2.On("Bounds").Return(&vec3.Box{})
	mockDb := new(db.MockObjects)
	mockDb.On("GetAll").Return(createGetManyResult(obj1, obj2))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	err := repo.loadFromDatabase()

	// Assert
	mockDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 2, rtree.Size())
}

func TestRepository_GetWithIDs_NoIDs_ReturnsEmpty(t *testing.T) {
	// Arrange
	mockDb := new(db.MockObjects)
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	results, err := flattenChannels(repo.GetWithIDs([]int64{}))

	// Assert
	assert.Empty(t, results)
	assert.NoError(t, err)
}

func TestRepository_GetWithIDs_DatabaseReturnsError_ReturnsError(t *testing.T) {
	// Arrange
	mockDb := new(db.MockObjects)
	mockDb.On("GetMany", []int64{1}).Return(createGetManyResult(errors.New("")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	results, err := flattenChannels(repo.GetWithIDs([]int64{1}))

	// Assert
	assert.Empty(t, results)
	assert.Error(t, err)
}

func TestRepository_GetWithIDs_TwoValidIDs_ReturnsTwoObjects(t *testing.T) {
	// Arrange
	obj1 := new(db.MockObject)
	obj1.On("ID").Return(int64(1))
	obj1.On("Bounds").Return(&vec3.Box{})
	obj2 := new(db.MockObject)
	obj2.On("ID").Return(int64(2))
	obj2.On("Bounds").Return(&vec3.Box{})
	mockDb := new(db.MockObjects)
	mockDb.On("GetMany", []int64{1, 2}).
		Return(createGetManyResult(obj1, obj2))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	results, err := flattenChannels(repo.GetWithIDs([]int64{1, 2}))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))
}

func TestRepository_GetWithID_InvalidID_ReturnsError(t *testing.T) {
	// Arrange
	mockDb := new(db.MockObjects)
	mockDb.On("GetMany", []int64{1}).
		Return(createGetManyResult(errors.New("")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	result, err := repo.GetWithID(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
func TestRepository_GetWithID_ValidID_ReturnsObject(t *testing.T) {
	// Arrange
	obj1 := new(db.MockObject)
	obj1.On("ID").Return(int64(1))
	mockDb := new(db.MockObjects)
	mockDb.On("GetMany", []int64{1}).
		Return(createGetManyResult(obj1))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	result, err := repo.GetWithID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
