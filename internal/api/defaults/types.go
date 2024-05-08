package defaults

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

type Defaults struct {
	spec         *openapi3.T
	typeMappings map[string]TypeMappings
	creators     map[string]Creator
}

type TypeMappings struct {
	Single func(isPointer bool, schema *openapi3.Schema, field reflect.Value)
	Slice  func(isPointer bool, schema *openapi3.Schema, field reflect.Value)
}

type Creator struct {
	New func() any
}
