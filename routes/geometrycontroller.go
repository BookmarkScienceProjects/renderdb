package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/larsmoa/renderdb/formats"
	"github.com/larsmoa/renderdb/repository"
	"github.com/larsmoa/renderdb/repository/options"

	"github.com/go-martini/martini"
	"github.com/ungerik/go3d/float64/vec3"
)

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

	repo repository.Repository
}

func (c *GeometryController) Init(repo repository.Repository, route *martini.ClassicMartini) {
	c.repo = repo
	route.Group("/geometry", func(r martini.Router) {
		r.Post("/", c.HandlePostGeometry)
		r.Post("/obj", c.HandlePostObjFile)
		r.Post("/view", c.HandlePostView)
		r.Get("/(?P<id>[0-9]+)", c.HandleGetGeometry)
	})
}

func (c *GeometryController) HandlePostGeometry(w http.ResponseWriter, r *http.Request, params martini.Params) {
	// Parse request
	ctx, httpErr := c.CreateContext(r, params)
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

func (c *GeometryController) HandlePostObjFile(w http.ResponseWriter, r *http.Request, params martini.Params) {
	// Parse body
	ctx, httpErr := c.CreateContext(r, params)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	defer ctx.Body.Close()
	reader := formats.WavefrontObjReader{}
	reader.SetOptions(formats.ReadOptions{DiscardDegeneratedFaces: true})
	err := reader.Read(ctx.Body)
	if err != nil {
		c.HandleError(w, NewHttpError(err, http.StatusBadRequest))
		return
	}
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
		obj := repository.NewSimpleObject(bounds, buf.Bytes(), nil)
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

func (c *GeometryController) HandlePostView(w http.ResponseWriter, r *http.Request, params martini.Params) {
	// Parse body
	ctx, httpErr := c.CreateContext(r, params)
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

	// Get object IDs inside view
	sortByDistanceOptions := options.SortByDistance{*payload.EyePosition}
	volume := *payload.Volume()
	ids, err := c.repo.GetInsideVolumeIDs(volume, sortByDistanceOptions)
	if err != nil {
		c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
		return
	}
	c.WriteResponse(w, ids)
}

func (c *GeometryController) HandleGetGeometry(w http.ResponseWriter, r *http.Request, params martini.Params) {
	ctx, httpErr := c.CreateContext(r, params)
	if httpErr != nil {
		c.HandleError(w, httpErr)
		return
	}

	var id int
	id, _ = strconv.Atoi(ctx.Vars["id"])
	object, err := c.repo.GetWithID(int64(id))
	if err != nil {
		c.HandleError(w, NewHttpError(err, http.StatusInternalServerError))
		return
	}
	c.WriteResponse(w, newGeometryObjectResponsePayload(object))
}

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

func (p *geometryRequestPayloadWrapper) ID() int64 {
	panic("Not supported")
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

type geometryObjectResponsePayload geometryRequestPayload

func newGeometryObjectResponsePayload(obj repository.Object) geometryObjectResponsePayload {
	payload := geometryObjectResponsePayload{}
	bounds := obj.Bounds()
	payload.Bounds = &BoundsPayload{bounds.Min, bounds.Max}
	payload.GeometryData = obj.GeometryData()
	buffer, _ := json.Marshal(obj.Metadata())
	payload.Metadata = string(buffer)
	return payload
}
