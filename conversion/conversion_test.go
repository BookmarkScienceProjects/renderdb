package conversion

import (
	"testing"

	"github.com/dhconnelly/rtreego"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestBoxToRect_MinIsOrigo_ReturnsCorrect(t *testing.T) {
	expected, _ := rtreego.NewRect(rtreego.Point{0, 0, 0}, rtreego.Point{1, 2, 3})
	rect := BoxToRect(&vec3.Box{vec3.T{0, 0, 0}, vec3.T{1, 2, 3}})
	assert.Equal(t, expected, rect)
}

func TestBoxRect_MinIsPostive_ReturnsCorrect(t *testing.T) {
	expected, _ := rtreego.NewRect(rtreego.Point{0.5, 1, 1.5}, rtreego.Point{1, 1, 1})
	rect := BoxToRect(&vec3.Box{vec3.T{0.5, 1, 1.5}, vec3.T{1.5, 2, 2.5}})
	assert.Equal(t, expected, rect)
}

func TestBoxRect_MinIsNegative_ReturnsCorrect(t *testing.T) {
	expected, _ := rtreego.NewRect(rtreego.Point{-1, -1, -1}, rtreego.Point{1, 1, 1})
	rect := BoxToRect(&vec3.Box{vec3.T{-1, -1, -1}, vec3.T{0, 0, 0}})
	assert.Equal(t, expected, rect)
}

func TestRectToBox_MinIsOrigo_ReturnsCorrect(t *testing.T) {
	rect, _ := rtreego.NewRect(rtreego.Point{0, 0, 0}, rtreego.Point{1, 2, 3})
	expected := &vec3.Box{vec3.T{0, 0, 0}, vec3.T{1, 2, 3}}
	assert.Equal(t, expected, RectToBox(rect))
}

func TestRectToBox_MinIsPostive_ReturnsCorrect(t *testing.T) {
	rect, _ := rtreego.NewRect(rtreego.Point{0.5, 1, 1.5}, rtreego.Point{1, 1, 1})
	expected := &vec3.Box{vec3.T{0.5, 1, 1.5}, vec3.T{1.5, 2, 2.5}}
	assert.Equal(t, expected, RectToBox(rect))
}

func TestRectToBox_MinIsNegative_ReturnsCorrect(t *testing.T) {
	rect, _ := rtreego.NewRect(rtreego.Point{-1, -1, -1}, rtreego.Point{1, 1, 1})
	expected := &vec3.Box{vec3.T{-1, -1, -1}, vec3.T{0, 0, 0}}
	assert.Equal(t, expected, RectToBox(rect))
}
