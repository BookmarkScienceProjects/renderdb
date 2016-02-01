package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestWavefrontObjLoader_ProcessMaterialLibrary_InvalidLine_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.Error(t, loader.processMaterialLibrary("invalid mtllib line"))
}

func TestWavefrontObjLoader_ProcessMaterialLibrary_ValidLine_SetsLibrary(t *testing.T) {
	loader := WavefrontObjLoader{}
	err := loader.processMaterialLibrary("mtllib      materials.mtl")
	assert.NoError(t, err)
	assert.Equal(t, "materials.mtl", loader.mtllib)
}

func TestWavefrontObjLoader_ProcessMaterialLibrary_AlreadySet_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	loader.mtllib = "somefile.mtl"
	assert.Error(t, loader.processMaterialLibrary("mtllib materials.mtl"))
}

func TestWavefrontObjLoader_ProcessGroup_ValidLine_EndsAndStartsFaceset(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}
	loader.g = append(loader.g, group{facesetCount: -1})

	// Act
	err := loader.processGroup("g   group")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, loader.g[0].facesetCount)
	assert.Equal(t, 2, len(loader.g))
	assert.Equal(t, "group", loader.g[1].name)
}

func TestWavefrontObjLoader_ProcessGroup_InvalidLine_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	err := loader.processUseMaterial("not a g line")
	assert.Error(t, err)
}

func TestWavefrontObjLoader_ProcessUseMaterial_ValidLine_EndsAndStartsFaceset(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}
	loader.facesets = append(loader.facesets, faceset{faceCount: -1})

	// Act
	err := loader.processUseMaterial("usemtl       material_name")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, loader.facesets[0].faceCount)
	assert.Equal(t, 2, len(loader.facesets))
	assert.Equal(t, "material_name", loader.facesets[1].material)
}

func TestWavefrontObjLoader_ProcessFace_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.Error(t, loader.processFace([]string{}))
	assert.Error(t, loader.processFace([]string{"a", "b", "c"}))
	assert.Error(t, loader.processFace([]string{"1/", "2/", "3/"}))
	assert.Error(t, loader.processFace([]string{"1/1", "2/2", "3/2"})) // Valid but not supported
	assert.Error(t, loader.processFace([]string{"1", "2"}))            // Too few coordinates
}

func TestWavefrontObjLoader_ProcessFace_VertexOnlyFormat_AddsFace(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	err := loader.processFace([]string{"1", "2", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.f))
	assert.Equal(t, 3, len(loader.f[0].corners))
	assert.Equal(t, 1, loader.f[0].corners[0].vertexIndex)
	assert.Equal(t, 2, loader.f[0].corners[1].vertexIndex)
	assert.Equal(t, 3, loader.f[0].corners[2].vertexIndex)
	assert.Equal(t, -1, loader.f[0].corners[0].normalIndex)
	assert.Equal(t, -1, loader.f[0].corners[1].normalIndex)
	assert.Equal(t, -1, loader.f[0].corners[2].normalIndex)
}

func TestWavefrontObjLoader_ProcessVertex_XYZ_AddsVertex(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	err := loader.processVertex([]string{"1.1", "2.0", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.v))
	assert.Equal(t, vec3.T{1.1, 2, 3}, loader.v[0])
}

func TestWavefrontObjLoader_ProcessVertex_XYZW_IgnoresW(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	err := loader.processVertex([]string{"1", "2", "3", "999"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.v))
	assert.Equal(t, vec3.T{1, 2, 3}, loader.v[0])
}

func TestWavefrontObjLoader_ProcessVertex_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.Error(t, loader.processVertex([]string{"0", "0"}))                // XY only
	assert.Error(t, loader.processVertex([]string{"0", "0", "A"}))           // Non-number
	assert.Error(t, loader.processVertex([]string{"0", "0", "0", "1", "2"})) // More than 4 coordinates
}

func TestWavefrontObjLoader_ProcessVertexNormal_XYZ_AddsNormal(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	err := loader.processVertexNormal([]string{"1.1", "2.0", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.vn))
	assert.Equal(t, vec3.T{1.1, 2, 3}, loader.vn[0])
}

func TestWavefrontObjLoader_ProcessVertexNormal_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.Error(t, loader.processVertexNormal([]string{"0", "0"}))           // XY only
	assert.Error(t, loader.processVertexNormal([]string{"0", "0", "A"}))      // Non-number
	assert.Error(t, loader.processVertexNormal([]string{"0", "0", "0", "1"})) // More than 3 coordinates
}

func TestWavefrontObjLoader_StartGroup_StartsNewGroup(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	loader.startGroup("MyGroup")

	// Assert
	assert.Equal(t, 1, len(loader.g))
	assert.Equal(t, "MyGroup", loader.g[0].name)
	assert.Equal(t, 0, loader.g[0].firstFacesetIndex)
	assert.Equal(t, -1, loader.g[0].facesetCount)
}

func TestWavefrontObjLoader_EndGroup_NoGroups_DoesNotPanic(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.NotPanics(t, func() {
		loader.endGroup()
	})
}

func TestWavefrontObjLoader_EndGroup_GroupStarted_UpdatesFacesetCount(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}
	loader.g = append(loader.g, group{
		name:              "Test",
		firstFacesetIndex: 0,
		facesetCount:      -1,
	})

	// Act
	loader.facesets = append(loader.facesets, faceset{})
	loader.endGroup()

	// Assert
	assert.Equal(t, 1, loader.g[0].facesetCount)
}

func TestWavefrontObjLoader_StartFaceset_StartsNewFaceset(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}

	// Act
	loader.startFaceset("SomeMaterial")

	// Assert
	assert.Equal(t, 1, len(loader.facesets))
	assert.Equal(t, "SomeMaterial", loader.facesets[0].material)
	assert.Equal(t, 0, loader.facesets[0].firstFaceIndex)
	assert.Equal(t, -1, loader.facesets[0].faceCount)
}

func TestWavefrontObjLoader_EndFaceset_NoFacesets_DoesNotPanic(t *testing.T) {
	loader := WavefrontObjLoader{}
	assert.NotPanics(t, func() {
		loader.endFaceset()
	})
}

func TestWavefrontObjLoader_EndFaceset_FacesetStarted_UpdatesFaceCount(t *testing.T) {
	// Arrange
	loader := WavefrontObjLoader{}
	loader.facesets = append(loader.facesets, faceset{
		material:       "Test",
		firstFaceIndex: 0,
		faceCount:      -1,
	})

	// Act
	loader.f = append(loader.f, face{})
	loader.endFaceset()

	// Assert
	assert.Equal(t, 1, loader.facesets[0].faceCount)
}
