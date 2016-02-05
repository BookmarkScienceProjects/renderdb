package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func createFace(cornerIdx ...int) face {
	f := face{}
	f.corners = make([]faceCorner, len(cornerIdx))
	for i := 0; i < len(cornerIdx); i++ {
		f.corners[i].vertexIndex = cornerIdx[i]
		f.corners[i].normalIndex = cornerIdx[i]
	}
	return f
}

func TestGroup_BuildFormats_EmptyGroup_ReturnsEmptyBuffer(t *testing.T) {
	// Arrange
	g := group{}
	origBuffer := objBuffer{}

	// Act
	buffer := g.buildBuffers(&origBuffer)

	// Assert
	assert.Equal(t, 0, len(buffer.facesets))
	assert.Equal(t, 0, len(buffer.f))
	assert.Equal(t, 0, len(buffer.v))
	assert.Equal(t, 0, len(buffer.vn))
}

func TestGroup_BuildFormats_SingleGroupWithSingleFace_ReturnsCorrect(t *testing.T) {
	// Arrange
	g := group{}
	g.firstFacesetIndex = 0
	g.facesetCount = 1

	fs := faceset{}
	fs.firstFaceIndex = 1
	fs.faceCount = 1
	fs.material = "Abc"

	origBuffer := objBuffer{}
	origBuffer.g = []group{g}
	origBuffer.f = []face{
		createFace(0, 1, 2),
	}
	origBuffer.facesets = []faceset{fs}
	origBuffer.v = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{1, 1, 1},
		vec3.T{2, 2, 2},
	}
	origBuffer.vn = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{-1, -1, -1},
		vec3.T{-2, -2, -2},
	}

	// Act
	buffer := g.buildBuffers(&origBuffer)

	// Assert
	assert.Equal(t, 1, len(buffer.g))
	assert.Equal(t, 1, len(buffer.facesets))
	assert.Equal(t, 1, len(buffer.f))
	assert.Equal(t, 3, len(buffer.v))
	assert.Equal(t, 3, len(buffer.vn))
}

func TestGroup_BuildFormats_GroupWithOneFaceset_ReturnsCorrectSubset(t *testing.T) {
	// Arrange
	origBuffer := objBuffer{}
	origBuffer.f = []face{
		// Faceset 1
		createFace(0, 2, 4),
		createFace(4, 2, 6),
		// Faceset 2
		createFace(1, 3, 5),
		createFace(5, 3, 7),
	}
	origBuffer.facesets = []faceset{
		faceset{
			firstFaceIndex: 0,
			faceCount:      2,
			material:       "Material 1",
		},
		faceset{
			firstFaceIndex: 2,
			faceCount:      2,
			material:       "Material 2",
		},
	}
	origBuffer.v = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{1, 1, 1},
		vec3.T{2, 2, 2},
		vec3.T{3, 3, 3},
		vec3.T{4, 4, 4},
		vec3.T{5, 5, 5},
		vec3.T{6, 6, 6},
		vec3.T{7, 7, 7},
	}
	origBuffer.vn = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{-1, -1, -1},
		vec3.T{-2, -2, -2},
		vec3.T{-3, -3, -3},
		vec3.T{-4, -4, -4},
		vec3.T{-5, -5, -5},
		vec3.T{-6, -6, -6},
		vec3.T{-7, -7, -7},
	}

	g1 := group{name: "Group 1", firstFacesetIndex: 0, facesetCount: 1}
	g2 := group{name: "Group 2", firstFacesetIndex: 1, facesetCount: 1}
	origBuffer.g = []group{g1, g2}

	// Act
	buffer := g1.buildBuffers(&origBuffer)

	// Assert
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{0, 0, 0}, vec3.T{2, 2, 2}, vec3.T{4, 4, 4}, vec3.T{6, 6, 6},
		},
		buffer.v)
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{0, 0, 0}, vec3.T{-2, -2, -2}, vec3.T{-4, -4, -4}, vec3.T{-6, -6, -6},
		},
		buffer.vn)
	assert.Equal(t, 1, len(buffer.facesets))
	assert.Equal(t,
		faceset{firstFaceIndex: 0, faceCount: 2, material: "Material 1"},
		buffer.facesets[0])
	assert.Equal(t, 1, len(buffer.g))
	assert.Equal(t,
		group{name: "Group 1", firstFacesetIndex: 0, facesetCount: 1},
		buffer.g[0])
}

func TestGroup_BuildFormats_GroupWithTwoFacesets_ReturnsCorrectSubset(t *testing.T) {
	// Arrange
	origBuffer := objBuffer{}
	origBuffer.f = []face{
		// Faceset 1
		createFace(0, 2, 4),
		createFace(4, 2, 6),
		// Faceset 2
		createFace(1, 3, 5),
		createFace(5, 3, 7),
		// Faceset 3
		createFace(5, 7, 2),
		createFace(7, 5, 4),
	}
	origBuffer.facesets = []faceset{
		faceset{
			firstFaceIndex: 0,
			faceCount:      2,
			material:       "Material 1",
		},
		faceset{
			firstFaceIndex: 2,
			faceCount:      2,
			material:       "Material 2",
		},
		faceset{
			firstFaceIndex: 4,
			faceCount:      2,
			material:       "Material 3",
		},
	}
	origBuffer.v = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{1, 1, 1},
		vec3.T{2, 2, 2},
		vec3.T{3, 3, 3},
		vec3.T{4, 4, 4},
		vec3.T{5, 5, 5},
		vec3.T{6, 6, 6},
		vec3.T{7, 7, 7},
	}
	origBuffer.vn = []vec3.T{
		vec3.T{0, 0, 0},
		vec3.T{-1, -1, -1},
		vec3.T{-2, -2, -2},
		vec3.T{-3, -3, -3},
		vec3.T{-4, -4, -4},
		vec3.T{-5, -5, -5},
		vec3.T{-6, -6, -6},
		vec3.T{-7, -7, -7},
	}

	g1 := group{name: "Group 1", firstFacesetIndex: 0, facesetCount: 2}
	g2 := group{name: "Group 2", firstFacesetIndex: 2, facesetCount: 1}
	origBuffer.g = []group{g1, g2}

	// Act
	buffer := g1.buildBuffers(&origBuffer)

	// Assert
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{0, 0, 0}, vec3.T{2, 2, 2}, vec3.T{4, 4, 4}, vec3.T{6, 6, 6},
		},
		buffer.v)
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{0, 0, 0}, vec3.T{-2, -2, -2}, vec3.T{-4, -4, -4}, vec3.T{-6, -6, -6},
		},
		buffer.vn)
	assert.Equal(t, 2, len(buffer.facesets))
	assert.Equal(t,
		faceset{firstFaceIndex: 0, faceCount: 2, material: "Material 1"},
		buffer.facesets[0])
	assert.Equal(t,
		faceset{firstFaceIndex: 2, faceCount: 2, material: "Material 2"},
		buffer.facesets[1])
	assert.Equal(t,
		group{firstFacesetIndex: 0, facesetCount: 2, name: "Group 1"},
		buffer.g[0])
}
