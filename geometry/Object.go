package geometry

import "github.com/ungerik/go3d/float64/vec3"

type Object interface {
	// The bound of the object
	Bounds() *vec3.Box

	GeometryData() []byte
	Metadata() interface{}
}

type SimpleObject struct {
	bounds       *vec3.Box
	geometryData []byte
	metadata     interface{}
}

func NewSimpleObject(bounds vec3.Box, geometryData []byte, metadata interface{}) *SimpleObject {
	return &SimpleObject{&bounds, geometryData, metadata}
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
