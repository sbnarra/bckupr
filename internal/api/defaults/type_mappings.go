package defaults

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

var basicTypeMappings = map[string]TypeMappings{
	"string": {
		Single: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSingleField(isPointer, field, schema.Default.(string))
		},
		Slice: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSliceField(isPointer, schema, field, func(i any) string {
				return i.(string)
			})
		},
	},
	"int": {
		Single: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSingleField(isPointer, field, int(schema.Default.(float64)))
		},
		Slice: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSliceField(isPointer, schema, field, func(i any) int {
				return i.(int)
			})
		},
	},
	"bool": {
		Single: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSingleField(isPointer, field, schema.Default.(bool))
		},
		Slice: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSliceField(isPointer, schema, field, func(i any) bool {
				return i.(bool)
			})
		},
	},
}

func (defaults *Defaults) getTypeMappings(fieldType string) (TypeMappings, error) {
	if typeMappings, found := defaults.typeMappings[fieldType]; found {
		return typeMappings, nil
	} else {
		return TypeMappings{}, errors.WithStack(fmt.Errorf("mapping missing: %v", fieldType))
	}
}

func setSliceField[T any](isPointer bool, schema *openapi3.Schema, field reflect.Value, convert func(i interface{}) T) {
	slice := []T{}
	for _, i := range schema.Default.([]interface{}) {
		slice = append(slice, convert(i))
	}
	setSingleField(isPointer, field, slice)
}

func setSingleField[T any](isPointer bool, field reflect.Value, val T) {
	if isPointer {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
}
