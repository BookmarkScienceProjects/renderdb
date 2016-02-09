package formats

import "github.com/larsmoa/renderdb/utils"

// group represents a named set of facesets.
type group struct {
	name              string
	firstFacesetIndex int
	facesetCount      int
}

func (g *group) buildBuffers(parentBuffer *objBuffer) *objBuffer {
	buffer := new(objBuffer)
	buffer.mtllib = parentBuffer.mtllib
	buffer.g = []group{
		group{
			name:         g.name,
			facesetCount: g.facesetCount,
		},
	}

	// Map from original vertex buffers to new
	vertexMapping := make([]int, len(parentBuffer.v))
	utils.FillIntSlice(vertexMapping, -1)
	normalMapping := make([]int, len(parentBuffer.vn))
	utils.FillIntSlice(normalMapping, -1)

	for i := g.firstFacesetIndex; i < g.firstFacesetIndex+g.facesetCount; i++ {
		originalFs := parentBuffer.facesets[i]
		// Create new faceset
		fs := faceset{}
		fs.firstFaceIndex = len(buffer.f)
		fs.faceCount = originalFs.faceCount
		fs.material = originalFs.material

		for j := fs.firstFaceIndex; j < fs.firstFaceIndex+fs.faceCount; j++ {
			originalFace := parentBuffer.f[j]

			// Create new face
			f := face{}
			f.corners = make([]faceCorner, len(originalFace.corners))

			for k, origCorner := range originalFace.corners {
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
				f.corners[k].vertexIndex, f.corners[k].normalIndex = newVertIdx, newNormIdx
			}

			buffer.f = append(buffer.f, f)
		}

		buffer.facesets = append(buffer.facesets, fs)
	}
	return buffer
}
