package geometry

import (
	"errors"
	"testing"

	"github.com/dhconnelly/rtreego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ungerik/go3d/vec3"
)

type mockGeometryDatabase struct {
	mock.Mock
}

func (m *mockGeometryDatabase) add(o GeometryObject) (int64, error) {
	args := m.Called(o)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockGeometryDatabase) getMany(ids []int64) (<-chan *geometryData, <-chan error) {
	args := m.Called(ids)
	return args.Get(0).(chan *geometryData), args.Get(1).(chan error)
}

// createGetManyResult is a helper function to emulate how getMany() returns
// values/error. Usage:
//   .Returns(createGetManyResult(data1, data2, err)) // Returns two values, then fails
func createGetManyResult(values ...interface{}) (chan *geometryData, chan error) {
	dataCh := make(chan *geometryData)
	errCh := make(chan error)
	go func() {
		defer close(dataCh)
		defer close(errCh)
		for _, v := range values {
			if data, ok := v.(*geometryData); ok {
				dataCh <- data
			} else if err, ok := v.(error); ok {
				errCh <- err
				return
			} else {
				panic("Values must be error or *geometryData")
			}
		}
	}()
	return dataCh, errCh
}

func flattenGetInsideVolumeResults(objCh <-chan GeometryObject, errCh <-chan error) ([]GeometryObject, error) {
	objects := []GeometryObject{}
	doBreak := false
	for !doBreak {
		select {
		case obj, ok := <-objCh:
			if ok {
				objects = append(objects, obj)
			} else {
				doBreak = true
			}
		case err, ok := <-errCh:
			if ok {
				return objects, err
			}
		}
	}
	return objects, nil
}

func TestAdd_ValidGeometry_AddsToTreeAndDatabase(t *testing.T) {
	// Arrange
	obj := new(SimpleGeometryObject)

	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(1, nil)

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

func TestAdd_DatabaseReturnsError_DoesNotAddToTree(t *testing.T) {
	// Arrange
	obj := new(SimpleGeometryObject)

	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(0, errors.New("error"))

	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}

	// Act
	_, err := repo.Add(obj)

	// Assert
	assert.Equal(t, 0, rtree.Size())
	assert.NotNil(t, err)
}

func TestGetInsideVolume_NothingInsideVolume_ReturnsEmpty(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{1, 1, 1}, vec3.T{2, 2, 2}}
	obj := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "",
		metadata:     nil,
	}

	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(1, nil)
	mockDb.On("getMany", []int64{}).Return(createGetManyResult())

	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	bounds := vec3.Box{vec3.T{5, 5, 5}, vec3.T{6, 6, 6}}
	objects, err := flattenGetInsideVolumeResults(repo.GetInsideVolume(bounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Empty(t, objects)
}

func TestGetInsideVolume_OneInsideVolume_ReturnsObject(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "",
		metadata:     nil,
	}

	data := new(geometryData)
	data.id = 1
	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(1, nil)
	mockDb.On("getMany", []int64{1}).Return(createGetManyResult(data))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := flattenGetInsideVolumeResults(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(objects))
}

func TestGetInsideVolume_DatabaseReturnsError_ReturnsError(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "",
		metadata:     nil,
	}

	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(1, nil)
	mockDb.On("getMany", []int64{1}).Return(createGetManyResult(errors.New("error")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := flattenGetInsideVolumeResults(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, 0, len(objects))
}

func TestGetInsideVolume_DatabaseReturnsOneThenError_ReturnsOneThenError(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj1 := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "1",
		metadata:     nil,
	}
	obj2 := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "2",
		metadata:     nil,
	}

	data := new(geometryData)
	data.id = 1
	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj1).Return(1, nil)
	mockDb.On("add", obj2).Return(2, nil)
	mockDb.On("getMany", []int64{1, 2}).Return(createGetManyResult(data, errors.New("error")))
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj1)
	repo.Add(obj2)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := flattenGetInsideVolumeResults(repo.GetInsideVolume(searchBounds))

	// Assert
	mockDb.AssertExpectations(t)
	assert.Equal(t, 1, len(objects))
	assert.Error(t, err)
}
