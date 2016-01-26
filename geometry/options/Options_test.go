package options

import (
	"testing"

	"github.com/dhconnelly/rtreego"
	"github.com/larsmoa/renderdb/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/vec3"
)

type stubFilterGeometryOption struct {
	keepIndices []int
}

func (s stubFilterGeometryOption) Apply(bounds []*vec3.Box) []int {
	return s.keepIndices
}

type stubSpatial struct {
	bounds vec3.Box
}

func (s stubSpatial) Bounds() *rtreego.Rect {
	return conversion.BoxToRect(&s.bounds)
}

func Test_VerifyAllAreOptions_Empty_ReturnsNoError(t *testing.T) {
	err := VerifyAllAreOptions()
	assert.NoError(t, err)
}

func Test_VerifyAllAreOptions_AllOptions_ReturnsNoError(t *testing.T) {
	// Arrange
	opt1 := new(stubFilterGeometryOption)
	opt2 := new(stubFilterGeometryOption)

	// Act
	err := VerifyAllAreOptions(opt1, opt2)

	// Assert
	assert.NoError(t, err)
}

func Test_VerifyAllAreOptions_NonOption_ReturnsError(t *testing.T) {
	err := VerifyAllAreOptions("somestring")
	assert.Error(t, err)
}

func Test_ApplyAllFilterGeometryOptions_NoOptions_ReturnsOriginal(t *testing.T) {
	// Arrange
	objects := []rtreego.Spatial{
		stubSpatial{vec3.Box{vec3.T{0, 0, 0}, vec3.T{1, 1, 1}}},
		stubSpatial{vec3.Box{vec3.T{0, 0, 0}, vec3.T{2, 2, 2}}},
	}

	// Act
	result := ApplyAllFilterGeometryOptions(objects)

	// Assert
	assert.EqualValues(t, objects, result)
}

func Test_ApplyAllFilterGeometryOptions_OptionKeepsFirst_ReturnsFirstOnly(t *testing.T) {
	// Arrange
	opt := stubFilterGeometryOption{[]int{0}}
	objects := []rtreego.Spatial{
		stubSpatial{vec3.Box{vec3.T{0, 0, 0}, vec3.T{1, 1, 1}}},
		stubSpatial{vec3.Box{vec3.T{0, 0, 0}, vec3.T{2, 2, 2}}},
	}

	// Act
	result := ApplyAllFilterGeometryOptions(objects, opt)

	// Assert
	assert.EqualValues(t, []rtreego.Spatial{objects[0]}, result)
}
