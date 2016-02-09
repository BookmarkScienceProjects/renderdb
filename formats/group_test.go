package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func createFace(material string, cornerIdx ...int) face {
	f := face{}
	f.corners = make([]faceCorner, len(cornerIdx))
	for i := 0; i < len(cornerIdx); i++ {
		f.corners[i].vertexIndex = cornerIdx[i]
		f.corners[i].normalIndex = cornerIdx[i]
	}
	f.material = material
	return f
}

func TestGroup_BuildFormats_EmptyGroup_ReturnsEmptyBuffer(t *testing.T) {
	// Arrange
	g := group{}
	origBuffer := objBuffer{}
	origBuffer.mtllib = "materials.mtl"

	// Act
	buffer := g.buildBuffers(&origBuffer)

	// Assert
	assert.Equal(t, "materials.mtl", buffer.mtllib)
	assert.Equal(t, 0, len(buffer.f))
	assert.Equal(t, 0, len(buffer.v))
	assert.Equal(t, 0, len(buffer.vn))
}

func TestGroup_BuildFormats_SingleGroupWithSingleFace_ReturnsCorrect(t *testing.T) {
	// Arrange
	g := group{}
	g.firstFaceIndex = 0
	g.faceCount = 1

	origBuffer := objBuffer{}
	origBuffer.g = []group{g}
	origBuffer.f = []face{
		createFace("mat", 0, 1, 2),
	}
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
	assert.Equal(t, 1, len(buffer.f))
	assert.Equal(t, 3, len(buffer.v))
	assert.Equal(t, 3, len(buffer.vn))
}

func TestGroup_BuildFormats_TwoGroupsWithTwoFaces_ReturnsCorrectGroups(t *testing.T) {
	// Arrange
	origBuffer := objBuffer{}
	origBuffer.f = []face{
		// Group 1
		createFace("mat1", 0, 2, 4),
		createFace("mat2", 4, 2, 6),
		// Group 2
		createFace("mat1", 1, 3, 5),
		createFace("mat2", 5, 3, 7),
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

	g1 := group{name: "Group 1", firstFaceIndex: 0, faceCount: 2}
	g2 := group{name: "Group 2", firstFaceIndex: 2, faceCount: 2}
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
	assert.Equal(t, 1, len(buffer.g))
	assert.Equal(t,
		group{name: "Group 1", firstFaceIndex: 0, faceCount: 2},
		buffer.g[0])
	assert.Equal(t, 2, len(buffer.f))
	assert.Equal(t, "mat1", buffer.f[0].material)
	assert.Equal(t, "mat2", buffer.f[1].material)
}

func TestGroup_BuildFormats_GroupWithTwoFacesets_ReturnsCorrectSubset(t *testing.T) {
	// Arrange
	origBuffer := objBuffer{}
	origBuffer.f = []face{
		// Group 1
		createFace("Material 1", 0, 2, 4),
		createFace("Material 1", 4, 2, 6),
		createFace("Material 2", 1, 3, 5),
		createFace("Material 2", 5, 3, 4),
		// Group 2
		createFace("Material 3", 5, 7, 2),
		createFace("Material 3", 7, 5, 4),
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

	g1 := group{name: "Group 1", firstFaceIndex: 0, faceCount: 4}
	g2 := group{name: "Group 2", firstFaceIndex: 4, faceCount: 2}
	origBuffer.g = []group{g1, g2}

	// Act
	buffer := g2.buildBuffers(&origBuffer)

	// Assert
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{5, 5, 5}, vec3.T{7, 7, 7}, vec3.T{2, 2, 2}, vec3.T{4, 4, 4},
		},
		buffer.v)
	assert.EqualValues(t,
		[]vec3.T{
			vec3.T{-5, -5, -5}, vec3.T{-7, -7, -7}, vec3.T{-2, -2, -2}, vec3.T{-4, -4, -4},
		},
		buffer.vn)
	assert.EqualValues(t, []face{
		createFace("Material 3", 0, 1, 2), // Remapped indices
		createFace("Material 3", 1, 0, 3), // Remapped indices
	}, buffer.f)
	assert.EqualValues(t, []group{group{"Group 2", 0, 2}}, buffer.g)
}
