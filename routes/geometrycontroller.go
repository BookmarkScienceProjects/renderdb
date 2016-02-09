package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/larsmoa/renderdb/formats"
	"github.com/larsmoa/renderdb/geometry"
	"github.com/larsmoa/renderdb/geometry/options"

	"github.com/gorilla/mux"
	"github.com/ungerik/go3d/float64/vec3"
)

type BoundsPayload struct {
	Min [3]float64 `json:"min"`
	Max [3]float64 `json:"max"`
}

type geometryRequestPayload struct {
	Bounds       *BoundsPayload `json:"bounds"`
	GeometryData []byte         `json:"geometryData"`
	Metadata     string         `json:"metadata"`
}

func (p *geometryRequestPayload) VerifyPayload() HttpError {
	if p.Bounds == nil {
		return NewHttpError(fmt.Errorf("Missing field 'bounds'"), http.StatusBadRequest)
	} else if p.GeometryData == nil {
		return NewHttpError(fmt.Errorf("Missing field 'geometryData'"), http.StatusBadRequest)
	} else if p.Metadata == "" {
		return NewHttpError(fmt.Errorf("Missing field 'metadata'"), http.StatusBadRequest)
	}
	return nil
}

type geometryRequestPayloadWrapper struct {
	payload *geometryRequestPayload
}

func (p *geometryRequestPayloadWrapper) Bounds() *vec3.Box {
	return &vec3.Box{p.payload.Bounds.Min, p.payload.Bounds.Max}
}

func (p *geometryRequestPayloadWrapper) GeometryData() []byte {
	return p.payload.GeometryData
}

func (p *geometryRequestPayloadWrapper) Metadata() interface{} {
	return p.payload.Metadata
}

type geometryResponsePayload struct {
	ID int64 `json:"id"`
}

type geometryViewRequestPayload struct {
	Bounds      *BoundsPayload `json:"bounds"`
	EyePosition *vec3.T        `json:"eyePosition"`
}

func (p *geometryViewRequestPayload) VerifyPayload() HttpError {
	if p.Bounds == nil {
		return NewHttpError(fmt.Errorf("Missing field 'bounds'"), http.StatusBadRequest)
	} else if p.EyePosition == nil {
		return NewHttpError(fmt.Errorf("Missing field 'eyePosition'"), http.StatusBadRequest)
	}
	return nil
}

func (p *geometryViewRequestPayload) Volume() *vec3.Box {
	return &vec3.Box{p.Bounds.Min, p.Bounds.Max}
}

type geometryViewResponsePayload geometryRequestPayload

func newViewResponsePayload(obj geometry.Object) geometryViewResponsePayload {
	payload := geometryViewResponsePayload{}
	bounds := obj.Bounds()
	payload.Bounds = &BoundsPayload{bounds.Min, bounds.Max}
	payload.GeometryData = obj.GeometryData()
	buffer, _ := json.Marshal(obj.Metadata())
	payload.Metadata = string(buffer)
	return payload
}

// GeometryController handles requests to "/geometry".
// Supported endpoints:
// - POST  /geometry
// -- Adds new geometry to the database
// -- Request body: {"bounds": {"min": [0,0,0], "max": [1,1,1]}, "geometryText": "...", "metadata": {...}}
// -- Response body: {"id": 1}
// - POST  /geometry/volume
// -- Requests objects within a volume
// -- Request body: {"bounds": {"min": [0,0,0], "max": [1,1,1]}, "eyePosition": [-1,-1,0]}
// -- Response body: [{"bounds": {"min": [0,0,0], "max": [1,1,1]}, "geometryText": "...", "metadata": {...}}, {...}, ..., {...}]
type GeometryController struct {
	Controller

	repo geometry.Repository
}

func (c *GeometryController) Init(repo geometry.Repository, route *mux.Router) {
	c.repo = repo
	route.Path("/geometry").Methods("POST").HandlerFunc(c.HandlePostGeometry)
	route.Path("/geometry/obj").Methods("POST").HandlerFunc(c.HandlePostObjFile)
	route.Path("/geometry/view").Methods("POST").HandlerFunc(c.HandlePostView)
}

func (c *GeometryController) HandlePostGeometry(w http.ResponseWriter, r *http.Request) {
	// Parse request
	ctx, httpErr := c.CreateContext(r)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}
	payload := new(geometryRequestPayload)
	httpErr = c.ParseBody(ctx, payload)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	// Verify payload
	httpErr = payload.VerifyPayload()
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	// Upload!
	id, err := c.repo.Add(&geometryRequestPayloadWrapper{payload})
	if err != nil {
		c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
		return
	}
	c.WriteResponse(w, geometryResponsePayload{id})
}

func (c *GeometryController) HandlePostView(w http.ResponseWriter, r *http.Request) {
	// Parse body
	ctx, httpErr := c.CreateContext(r)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}
	payload := new(geometryViewRequestPayload)
	httpErr = c.ParseBody(ctx, payload)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	// Verify payload
	httpErr = payload.VerifyPayload()
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	// Get objects inside view
	sortByDistanceOptions := options.SortByDistance{*payload.EyePosition}
	objects := make([]geometryViewResponsePayload, 0, 100)
	volume := *payload.Volume()
	objCh, errCh := c.repo.GetInsideVolume(volume, sortByDistanceOptions)
	more := true
	for more {
		var obj geometry.Object
		var err error
		select {
		case obj, more = <-objCh:
			if more {
				objects = append(objects, newViewResponsePayload(obj))
			}

		case err = <-errCh:
			c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
			return
		}
	}

	c.WriteResponse(w, objects)
}

func (c *GeometryController) HandlePostObjFile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	ctx, httpErr := c.CreateContext(r)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	defer ctx.Body.Close()
	reader := formats.WavefrontObjReader{}
	err := reader.Read(ctx.Body)
	if err != nil {
		c.HandleError(w, NewHttpError(err, http.StatusBadRequest))
		return
	}

	f, _ := os.Create("/tmp/test.obj")
	reader.Write(f)

	storedObjects := make(map[int64]int)
	groupCh := reader.Groups()
	for g := range groupCh {
		bounds := g.BoundingBox()
		buf := bytes.Buffer{}
		err = g.Write(&buf)
		if err != nil {
			c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
			return
		}

		var id int64
		obj := geometry.NewSimpleObject(bounds, buf.Bytes(), nil)
		id, err = c.repo.Add(obj)
		if err != nil {
			c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
			return
		}
		storedObjects[id] = buf.Len()
	}

	bytes, _ := json.Marshal(storedObjects)
	w.Write(bytes)
}
