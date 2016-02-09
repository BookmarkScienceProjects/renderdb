package formats

import "github.com/larsmoa/renderdb/utils"

// group represents a named set of facesets.
type group struct {
	name           string
	firstFaceIndex int
	faceCount      int
}

func (g *group) buildBuffers(parentBuffer *objBuffer) *objBuffer {
	buffer := new(objBuffer)
	buffer.mtllib = parentBuffer.mtllib
	buffer.g = []group{
		group{
			name:      g.name,
			faceCount: g.faceCount,
		},
	}

	// Map from original vertex buffers to new
	vertexMapping := make([]int, len(parentBuffer.v))
	utils.FillIntSlice(vertexMapping, -1)
	normalMapping := make([]int, len(parentBuffer.vn))
	utils.FillIntSlice(normalMapping, -1)

	for i := g.firstFaceIndex; i < g.firstFaceIndex+g.faceCount; i++ {

		originalFace := parentBuffer.f[i]

		// Create new face
		f := face{material: originalFace.material}
		f.corners = make([]faceCorner, len(originalFace.corners))

		for j, origCorner := range originalFace.corners {
			// Create new 'corners' and map indices
			origVertIdx := origCorner.vertexIndex
			origNormIdx := origCorner.normalIndex

			// Lookup or add new vertex
			var newVertIdx int
			if newVertIdx = vertexMapping[origVertIdx]; newVertIdx == -1 {
				newVertIdx = len(buffer.v)
				buffer.v = append(buffer.v, parentBuffer.v[origVertIdx])
				vertexMapping[origVertIdx] = newVertIdx
			}
			// Lookup or add new normal
			var newNormIdx int
			if newNormIdx = normalMapping[origNormIdx]; newNormIdx == -1 {
				newNormIdx = len(buffer.vn)
				buffer.vn = append(buffer.vn, parentBuffer.vn[origNormIdx])
				normalMapping[origNormIdx] = newNormIdx
			}

			// Add face corner
			f.corners[j].vertexIndex, f.corners[j].normalIndex = newVertIdx, newNormIdx
		}

		buffer.f = append(buffer.f, f)
	}
	return buffer
}
