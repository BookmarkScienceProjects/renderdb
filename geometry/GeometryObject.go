package geometry

import "github.com/ungerik/go3d/vec3"

type GeometryObject interface {
	// The bound of the object
	Bounds() vec3.Box

	GeometryText() string
	Metadata() interface{}
}

type SimpleGeometryObject struct {
	bounds       vec3.Box
	geometryText string
	metadata     interface{}
}

func (o *SimpleGeometryObject) Bounds() vec3.Box {
	return o.bounds
}

func (o *SimpleGeometryObject) GeometryText() string {
	return o.geometryText
}

func (o *SimpleGeometryObject) Metadata() interface{} {
	return o.metadata
}
