package formats

import (
	"fmt"
	"io"

	"github.com/ungerik/go3d/float64/vec3"
)

// Write outputs the buffer to the writer given. Returns an error if the
// operation fails.
func (b *objBuffer) Write(w io.Writer) error {
	var err error
	_, err = io.WriteString(w,
		fmt.Sprintf("# Exported using RenderDB\n"+
			"# %d vertices, %d normals, %d faces (%d facesets)\n",
			len(b.v), len(b.vn), len(b.f), len(b.facesets)))
	if err != nil {
		return err
	}
	if b.mtllib != "" {
		_, err = io.WriteString(w, fmt.Sprintf("mtllib %s\n", b.mtllib))
		if err != nil {
			return err
		}
	}
	if err = b.writeVertices(w); err != nil {
		return err
	}
	if err = b.writeNormals(w); err != nil {
		return err
	}
	for _, g := range b.g {
		if err = b.writeGroup(w, g); err != nil {
			return err
		}
	}

	return nil
}

func (b *objBuffer) writeVertices(w io.Writer) error {
	return writeVectors(w, "v %g %g %g\n", b.v)
}

func (b *objBuffer) writeNormals(w io.Writer) error {
	return writeVectors(w, "vn %g %g %g\n", b.vn)
}

func writeFace(w io.Writer, f face) error {
	var err error

	_, err = io.WriteString(w, "f")
	if err != nil {
		return err
	}

	for _, c := range f.corners {
		if c.normalIndex != -1 {
			_, err = io.WriteString(w,
				fmt.Sprintf(" %d//%d", c.vertexIndex+1, c.normalIndex+1))
		} else {
			_, err = io.WriteString(w, fmt.Sprintf(" %d", c.vertexIndex+1))
		}
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, "\n")
	return err
}

func writeVectors(w io.Writer, format string, vectors []vec3.T) error {
	for _, v := range vectors {
		_, err := io.WriteString(w, fmt.Sprintf(format, v[0], v[1], v[2]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *objBuffer) writeGroup(w io.Writer, g group) error {
	var err error
	_, err = io.WriteString(w, fmt.Sprintf("g %s", g.name))
	if err != nil {
		return err
	}
	for i := g.firstFacesetIndex; i < g.firstFacesetIndex+g.facesetCount; i++ {
		fs := b.facesets[i]
		if err = b.writeFaceset(w, fs); err != nil {
			return err
		}
	}

	return nil
}

func (b *objBuffer) writeFaceset(w io.Writer, fs faceset) error {
	var err error
	if fs.material != "" {
		_, err = io.WriteString(w, fmt.Sprintf("usemtl %s\n", fs.material))
		if err != nil {
			return err
		}
	}

	for i := fs.firstFaceIndex; i < fs.firstFaceIndex+fs.faceCount; i++ {
		if err = writeFace(w, b.f[i]); err != nil {
			return err
		}
	}

	return nil
}
