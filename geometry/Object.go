package geometry

import "github.com/ungerik/go3d/vec3"

type Object interface {
	// The bound of the object
	Bounds() *vec3.Box

	GeometryText() string
	Metadata() interface{}
}

type SimpleObject struct {
	bounds       *vec3.Box
	geometryText string
	metadata     interface{}
}

func (o *SimpleObject) Bounds() *vec3.Box {
	return o.bounds
}

func (o *SimpleObject) GeometryText() string {
	return o.geometryText
}

func (o *SimpleObject) Metadata() interface{} {
	return o.metadata
}
