package db

import "github.com/ungerik/go3d/float64/vec3"

type Object interface {
	// ID returns an unique ID of the object
	ID() int64

	// WorldID returns the ID of the world the object 'lives' in
	WorldID() int64
	// LayerID returns the ID of the layer the object is a part of
	LayerID() int64
	// SceneID returns the ID of the scene the object is a part of
	SceneID() int64

	// Bounds returns the bounding box of the object.
	Bounds() *vec3.Box
	// GeometryData returns raw geometry data of the object
	GeometryData() []byte
	// Metadata returns arbitrary JSON-convertible metadata for
	// the object.
	Metadata() interface{}
}

type SimpleObject struct {
	id           int64
	worldID      int64
	layerID      int64
	sceneID      int64
	bounds       *vec3.Box
	geometryData []byte
	metadata     interface{}
}

func NewSimpleObject(bounds vec3.Box, geometryData []byte, metadata interface{}) *SimpleObject {
	o := new(SimpleObject)
	o.id = -1
	o.worldID = -1
	o.layerID = -1
	o.sceneID = -1
	o.bounds = &bounds
	o.geometryData = geometryData
	o.metadata = metadata
	return o
}

func (o *SimpleObject) ID() int64 {
	return o.id
}

func (o *SimpleObject) WorldID() int64 {
	return o.worldID
}
func (o *SimpleObject) LayerID() int64 {
	return o.layerID
}
func (o *SimpleObject) SceneID() int64 {
	return o.sceneID
}

func (o *SimpleObject) Bounds() *vec3.Box {
	return o.bounds
}

func (o *SimpleObject) GeometryData() []byte {
	return o.geometryData
}

func (o *SimpleObject) Metadata() interface{} {
	return o.metadata
}
