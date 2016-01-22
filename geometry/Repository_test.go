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

func (m *mockGeometryDatabase) getMany(ids []int64) ([]*geometryData, error) {
	args := m.Called(ids)
	return args.Get(0).([]*geometryData), args.Error(1)
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
	mockDb.On("getMany", []int64{}).Return([]*geometryData{}, nil)
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	bounds := vec3.Box{vec3.T{5, 5, 5}, vec3.T{6, 6, 6}}
	objects, err := repo.GetInsideVolume(bounds)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, objects)
	mockDb.AssertExpectations(t)
}

func TestGetInsideVolume_OneInsideVolume_ReturnsObject(t *testing.T) {
	// Arrange
	objBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	obj := &SimpleGeometryObject{
		bounds:       objBounds,
		geometryText: "",
		metadata:     nil,
	}

	mockDb := new(mockGeometryDatabase)
	mockDb.On("add", obj).Return(1, nil)
	mockDb.On("getMany", []int64{1}).Return([]*geometryData{}, nil)
	rtree := rtreego.NewTree(3, 5, 10)
	repo := defaultRepository{mockDb, rtree}
	repo.Add(obj)

	// Act
	searchBounds := vec3.Box{vec3.T{0.5, 0.5, 0.5}, vec3.T{1.5, 1.5, 1.5}}
	objects, err := repo.GetInsideVolume(searchBounds)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(objects))
	mockDb.AssertExpectations(t)
}
