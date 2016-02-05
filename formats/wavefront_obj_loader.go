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
	groupRegex = regexp.MustCompile(`^g\s+(.*)$`)
	usemtlRegex = regexp.MustCompile(`^usemtl\s+(.*)$`)
	mtllibRegex = regexp.MustCompile(`^mtllib\s+(.*)$`)
}

type lineError struct {
	lineNumber int
	line       string
	err        error
}

func (e lineError) Error() string {
	return fmt.Sprintf("Line #%d: %v ('%s')", e.lineNumber, e.line, e.err)
}

// faceCorner represents a 'corner' (or vertex) in a face
type faceCorner struct {
	vertexIndex int
	normalIndex int
}

// face represents a surface represented by a set of corner
type face struct {
	corners []faceCorner
}

// faceset represents a set of faces that share the same
// material.
type faceset struct {
	firstFaceIndex int
	faceCount      int
	material       string
}

type WavefrontObjLoader struct {
	objBuffer
}

func (l *WavefrontObjLoader) Load(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		i++
		if len(line) == 0 {
			continue
		}

		if !strings.HasPrefix(line, "#") { // Ignore comments
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
			}

			if err != nil {
				return lineError{i, line, err}
			}
		}
	}

	return scanner.Err()
}

func (l *WavefrontObjLoader) processVertex(fields []string) error {
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

func (l *WavefrontObjLoader) processVertexNormal(fields []string) error {
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
		return faceCorner{v, -1}, err
	} else if match := faceVertexAndNormalRegex.FindStringSubmatch(field); match != nil {
		// f v1//n1 v2//n2 ... - vertex and normal
		v, errV := strconv.Atoi(match[1])
		n, errN := strconv.Atoi(match[2])
		return faceCorner{v, n}, utils.FirstError(errV, errN)
	} else {
		// Note! f v1/t1 v2/t2 ... - vertex + texture and
		// f v1/t1/n1 v2/t2/n2 ... - vertex, texture and normal
		// are not currently supported.
		return faceCorner{-1, -1}, fmt.Errorf("Face field '%s' is not on a supported format", field)
	}
}

func (l *WavefrontObjLoader) processFace(fields []string) error {
	if len(fields) < 3 {
		return fmt.Errorf("Expected %d fields, but got %d", 3, len(fields))
	}

	f := face{make([]faceCorner, len(fields))}
	for i, field := range fields {
		corner, err := parseFaceField(field)
		if err != nil {
			return err
		}
		f.corners[i] = corner
	}
	l.f = append(l.f, f)
	return nil
}

func (l *WavefrontObjLoader) processGroup(line string) error {
	if match := groupRegex.FindStringSubmatch(line); match != nil {
		l.endGroup()
		l.startGroup(match[1])
		return nil
	}
	return fmt.Errorf("Could not parse group")
}

func (l *WavefrontObjLoader) processMaterialLibrary(line string) error {
	if l.mtllib != "" {
		return fmt.Errorf("Material library already set")
	}
	if match := mtllibRegex.FindStringSubmatch(line); match != nil {
		l.mtllib = match[1]
		return nil
	}
	return fmt.Errorf("Could not parse 'mtllib'-line")
}

func (l *WavefrontObjLoader) processUseMaterial(line string) error {
	if match := usemtlRegex.FindStringSubmatch(line); match != nil {
		l.endFaceset()
		l.startFaceset(match[1])
		return nil
	}
	return fmt.Errorf("Could not parse 'usemtl'-line")
}

func (l *WavefrontObjLoader) startFaceset(material string) {
	fs := faceset{
		material:       material,
		firstFaceIndex: len(l.f),
		faceCount:      -1,
	}
	l.facesets = append(l.facesets, fs)
}

func (l *WavefrontObjLoader) endFaceset() {
	if len(l.facesets) > 0 {
		lastIdx := len(l.facesets) - 1
		l.facesets[lastIdx].faceCount = len(l.f) - l.facesets[lastIdx].firstFaceIndex
	}
}

func (l *WavefrontObjLoader) startGroup(name string) {
	g := group{
		name:              name,
		firstFacesetIndex: len(l.facesets),
		facesetCount:      -1,
	}
	l.g = append(l.g, g)
}
func (l *WavefrontObjLoader) endGroup() {
	if len(l.g) > 0 {
		idx := len(l.g) - 1
		l.g[idx].facesetCount = len(l.facesets) - l.g[idx].firstFacesetIndex
	}
}
