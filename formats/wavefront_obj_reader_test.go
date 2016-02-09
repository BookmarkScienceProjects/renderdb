package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestWavefrontObjReader_ProcessMaterialLibrary_InvalidLine_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	assert.Error(t, loader.processMaterialLibrary("invalid mtllib line"))
}

func TestWavefrontObjReader_ProcessMaterialLibrary_ValidLine_SetsLibrary(t *testing.T) {
	loader := WavefrontObjReader{}
	err := loader.processMaterialLibrary("mtllib      materials.mtl")
	assert.NoError(t, err)
	assert.Equal(t, "materials.mtl", loader.mtllib)
}

func TestWavefrontObjReader_ProcessMaterialLibrary_AlreadySet_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	loader.mtllib = "somefile.mtl"
	assert.Error(t, loader.processMaterialLibrary("mtllib materials.mtl"))
}

func TestWavefrontObjReader_ProcessGroup_ValidLine_EndsAndStartsGroup(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}
	loader.f = []face{face{}}
	loader.g = append(loader.g, group{firstFaceIndex: 0, faceCount: -1})

	// Act
	err := loader.processGroup("g   group")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, loader.g[0].faceCount)
	assert.Equal(t, 2, len(loader.g))
	assert.Equal(t, "group", loader.g[1].name)
}

func TestWavefrontObjReader_ProcessGroup_InvalidLine_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	err := loader.processUseMaterial("not a g line")
	assert.Error(t, err)
}

func TestWavefrontObjReader_ProcessUseMaterial_ValidLine_SetsActiveMaterial(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}
	loader.f = []face{face{}}

	// Act
	err := loader.processUseMaterial("usemtl       material_name")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "material_name", loader.activeMaterial)
}

func TestWavefrontObjReader_ProcessFace_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	assert.Error(t, loader.processFace([]string{}))
	assert.Error(t, loader.processFace([]string{"a", "b", "c"}))
	assert.Error(t, loader.processFace([]string{"1/", "2/", "3/"}))
	assert.Error(t, loader.processFace([]string{"1/1", "2/2", "3/2"})) // Valid but not supported
	assert.Error(t, loader.processFace([]string{"1", "2"}))            // Too few coordinates
}

func TestWavefrontObjReader_ProcessFace_VertexOnlyFormat_AddsFace(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	err := loader.processFace([]string{"1", "2", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.f))
	assert.Equal(t, 3, len(loader.f[0].corners))
	// Zero-based indices
	assert.Equal(t, 0, loader.f[0].corners[0].vertexIndex)
	assert.Equal(t, 1, loader.f[0].corners[1].vertexIndex)
	assert.Equal(t, 2, loader.f[0].corners[2].vertexIndex)
	assert.Equal(t, -1, loader.f[0].corners[0].normalIndex)
	assert.Equal(t, -1, loader.f[0].corners[1].normalIndex)
	assert.Equal(t, -1, loader.f[0].corners[2].normalIndex)
}

func TestWavefrontObjReader_ProcessVertex_XYZ_AddsVertex(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	err := loader.processVertex([]string{"1.1", "2.0", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.v))
	assert.Equal(t, vec3.T{1.1, 2, 3}, loader.v[0])
}

func TestWavefrontObjReader_ProcessVertex_XYZW_IgnoresW(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	err := loader.processVertex([]string{"1", "2", "3", "999"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.v))
	assert.Equal(t, vec3.T{1, 2, 3}, loader.v[0])
}

func TestWavefrontObjReader_ProcessVertex_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	assert.Error(t, loader.processVertex([]string{"0", "0"}))                // XY only
	assert.Error(t, loader.processVertex([]string{"0", "0", "A"}))           // Non-number
	assert.Error(t, loader.processVertex([]string{"0", "0", "0", "1", "2"})) // More than 4 coordinates
}

func TestWavefrontObjReader_ProcessVertexNormal_XYZ_AddsNormal(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	err := loader.processVertexNormal([]string{"1.1", "2.0", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.vn))
	assert.Equal(t, vec3.T{1.1, 2, 3}, loader.vn[0])
}

func TestWavefrontObjReader_ProcessVertexNormal_InvalidFields_ReturnsError(t *testing.T) {
	loader := WavefrontObjReader{}
	assert.Error(t, loader.processVertexNormal([]string{"0", "0"}))           // XY only
	assert.Error(t, loader.processVertexNormal([]string{"0", "0", "A"}))      // Non-number
	assert.Error(t, loader.processVertexNormal([]string{"0", "0", "0", "1"})) // More than 3 coordinates
}

func TestWavefrontObjReader_StartGroup_StartsNewGroup(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	loader.startGroup("MyGroup")

	// Assert
	assert.Equal(t, 1, len(loader.g))
	assert.Equal(t, "MyGroup", loader.g[0].name)
	assert.Equal(t, 0, loader.g[0].firstFaceIndex)
	assert.Equal(t, -1, loader.g[0].faceCount)
}

func TestWavefrontObjReader_EndGroup_NoGroups_DoesNotPanic(t *testing.T) {
	loader := WavefrontObjReader{}
	assert.NotPanics(t, func() {
		loader.endGroup()
	})
}

func TestWavefrontObjReader_EndGroup_GroupStarted_UpdatesFaceCount(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}
	loader.g = append(loader.g, group{
		name:           "Test",
		firstFaceIndex: 0,
		faceCount:      -1,
	})

	// Act
	loader.f = append(loader.f, createFace("mat", 1, 2, 3))
	loader.endGroup()

	// Assert
	assert.Equal(t, []group{group{"Test", 0, 1}}, loader.g)
}

func TestWavefrontObjReader_ProcessFace_UsesActiveMaterial(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}
	loader.activeMaterial = "my-material"

	// Act
	err := loader.processFace([]string{"1", "2", "3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loader.f))
	assert.Equal(t, "my-material", loader.f[0].material)
}

func TestWavefrontObjReader_EndGroup_EmptyGroup_IsDiscarded(t *testing.T) {
	// Arrange
	loader := WavefrontObjReader{}

	// Act
	loader.startGroup("someGroup")
	loader.endGroup()

	// Assert
	assert.Empty(t, loader.f)
	assert.Empty(t, loader.g)
}
