package formats

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/larsmoa/renderdb/utils"
	"github.com/ungerik/go3d/float64/vec3"
)

var faceVertexOnlyRegex *regexp.Regexp
var faceVertexAndNormalRegex *regexp.Regexp
var groupRegex *regexp.Regexp
var usemtlRegex *regexp.Regexp
var mtllibRegex *regexp.Regexp

func init() {
	faceVertexOnlyRegex = regexp.MustCompile(`^(\d+)$`)
	faceVertexAndNormalRegex = regexp.MustCompile(`^(\d+)//(\d+)$`)
	groupRegex = regexp.MustCompile(`^g\s*(.*)$`)
	usemtlRegex = regexp.MustCompile(`^usemtl\s+(.*)$`)
	mtllibRegex = regexp.MustCompile(`^mtllib\s+(.*)$`)
}

// WavefrontObjReader reads Wavefront OBJ files. The reader supports the
// following keywords:
// - v
// - vn
// - f
// - g
// - mtllib
// - usemtl
// The following keywords are ignored:
// - o
// - s
// - vt
// - vp
// - Comments (#)
// The reader supports splitting the OBJ file into 'groups' defined by the
// 'g'-keyword.
type WavefrontObjReader struct {
	objBuffer

	options ReadOptions
}

// SetOptions sets the read options that alters the behavior of
// Read. Defaults to the default ReadOptions{} struct.
func (l *WavefrontObjReader) SetOptions(options ReadOptions) {
	l.options = options
}

func (l *WavefrontObjReader) Read(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		i++
		// Ignore comments
		if hashPos := strings.IndexRune(line, '#'); hashPos != -1 {
			line = line[0:hashPos]
		}
		if len(line) == 0 {
			continue
		}

		var err error
		fields := strings.Fields(line)
		switch strings.ToLower(fields[0]) {
		case "v":
			err = l.processVertex(fields[1:])
		case "vn":
			err = l.processVertexNormal(fields[1:])
		case "f":
			err = l.processFace(fields[1:])
		case "g":
			err = l.processGroup(line)
		case "mtllib":
			err = l.processMaterialLibrary(line)
		case "usemtl":
			err = l.processUseMaterial(line)

			// Ignored keywords
		case "o":
		case "s":
		case "vt":
		case "vp":
			break

		default:
			err = fmt.Errorf("Unknown keyword '%s'", fields[0])
		}

		if err != nil {
			return lineError{i, line, err}
		}
	}
	l.endGroup()
	return scanner.Err()
}

// Groups returns a buffered channel with one element for each
// group in the loaded OBJ file.
func (l *WavefrontObjReader) Groups() <-chan GeometryGroup {
	ch := make(chan GeometryGroup, 10)
	go func() {
		defer close(ch)
		for _, g := range l.g {
			adapter := createGeometryGroupAdapter(&l.objBuffer, g)
			ch <- adapter
		}
	}()
	return ch
}

func (l *WavefrontObjReader) processVertex(fields []string) error {
	if len(fields) != 3 && len(fields) != 4 {
		return fmt.Errorf("Expected 3 or 4 fields, but got %d", len(fields))
	}
	x, errX := strconv.ParseFloat(fields[0], 64)
	y, errY := strconv.ParseFloat(fields[1], 64)
	z, errZ := strconv.ParseFloat(fields[2], 64)
	if err := utils.FirstError(errX, errY, errZ); err != nil {
		return err
	}
	l.v = append(l.v, vec3.T{x, y, z})
	return nil
}

func (l *WavefrontObjReader) processVertexNormal(fields []string) error {
	if len(fields) != 3 {
		return fmt.Errorf("Expected 3 fields, but got %d", len(fields))
	}
	x, errX := strconv.ParseFloat(fields[0], 64)
	y, errY := strconv.ParseFloat(fields[1], 64)
	z, errZ := strconv.ParseFloat(fields[2], 64)
	if err := utils.FirstError(errX, errY, errZ); err != nil {
		return err
	}
	l.vn = append(l.vn, vec3.T{x, y, z})
	return nil
}

func parseFaceField(field string) (faceCorner, error) {
	if match := faceVertexOnlyRegex.FindStringSubmatch(field); match != nil {
		// f v1 v2 ... - only vertex
		v, err := strconv.Atoi(match[1])
		return faceCorner{v - 1, -1}, err
	} else if match := faceVertexAndNormalRegex.FindStringSubmatch(field); match != nil {
		// f v1//n1 v2//n2 ... - vertex and normal
		v, errV := strconv.Atoi(match[1])
		n, errN := strconv.Atoi(match[2])
		return faceCorner{v - 1, n - 1}, utils.FirstError(errV, errN)
	} else {
		// Note! f v1/t1 v2/t2 ... - vertex + texture and
		// f v1/t1/n1 v2/t2/n2 ... - vertex, texture and normal
		// are not currently supported.
		return faceCorner{-1, -1}, fmt.Errorf("Face field '%s' is not on a supported format", field)
	}
}

func (l *WavefrontObjReader) isFaceAccepted(f *face) bool {
	if l.options.DiscardDegeneratedFaces {
		// Degenerated when a vertex index appears twice
		occurences := make(map[int]bool, len(f.corners))
		for _, c := range f.corners {
			vIdx := c.vertexIndex
			if _, ok := occurences[vIdx]; ok {
				return false // vIdx occurs twice, degenerated
			}
			occurences[vIdx] = true
		}
	}
	return true
}

func (l *WavefrontObjReader) processFace(fields []string) error {
	if len(fields) < 3 {
		return fmt.Errorf("Expected %d fields, but got %d", 3, len(fields))
	}

	f := face{make([]faceCorner, len(fields)), l.activeMaterial}
	for i, field := range fields {
		corner, err := parseFaceField(field)
		if err != nil {
			return err
		}
		f.corners[i] = corner
	}
	if l.isFaceAccepted(&f) {
		l.f = append(l.f, f)
	}
	return nil
}

func (l *WavefrontObjReader) processGroup(line string) error {
	if match := groupRegex.FindStringSubmatch(line); match != nil {
		l.endGroup()
		l.startGroup(match[1])
		return nil
	}
	return fmt.Errorf("Could not parse group")
}

func (l *WavefrontObjReader) processMaterialLibrary(line string) error {
	if l.mtllib != "" {
		return fmt.Errorf("Material library already set")
	}
	if match := mtllibRegex.FindStringSubmatch(line); match != nil {
		l.mtllib = match[1]
		return nil
	}
	return fmt.Errorf("Could not parse 'mtllib'-line")
}

func (l *WavefrontObjReader) processUseMaterial(line string) error {
	if match := usemtlRegex.FindStringSubmatch(line); match != nil {
		l.activeMaterial = match[1]
		return nil
	}
	return fmt.Errorf("Could not parse 'usemtl'-line")
}

func (l *WavefrontObjReader) startGroup(name string) {
	g := group{
		name:           name,
		firstFaceIndex: len(l.f),
		faceCount:      -1,
	}
	l.g = append(l.g, g)
}

func (l *WavefrontObjReader) isGroupAccepted(f *face) bool {
	if l.options.DiscardDegeneratedFaces {
		// Degenerated when a vertex index appears twice
		occurences := make(map[int]bool, len(f.corners))
		for _, c := range f.corners {
			vIdx := c.vertexIndex
			if _, ok := occurences[vIdx]; ok {
				return false // vIdx occurs twice, degenerated
			}
			occurences[vIdx] = true
		}
	}
	return true
}

func (l *WavefrontObjReader) endGroup() {
	if len(l.g) > 0 {
		idx := len(l.g) - 1
		count := len(l.f) - l.g[idx].firstFaceIndex
		if count > 0 {
			l.g[idx].faceCount = count
		} else {
			// Empty group, discard
			if len(l.g) > 0 {
				l.g = l.g[:len(l.g)-1]
			} else {
				l.g = nil
			}
		}
	}
}
