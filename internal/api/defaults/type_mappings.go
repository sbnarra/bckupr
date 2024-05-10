package defaults

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

var basicTypeMappings = map[string]TypeMappings{
	"string": standardTypeMapping(func(i any) string {
		return i.(string)
	}),
	"int": standardTypeMapping(func(i any) int {
		return i.(int)
	}),
	"bool": standardTypeMapping(func(i any) bool {
		return i.(bool)
	}),
}

func standardTypeMapping[O any](t func(any) O) TypeMappings {
	return TypeMappings{
		Single: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSingleField(isPointer, field, t(schema.Default))
		},
		Slice: func(isPointer bool, schema *openapi3.Schema, field reflect.Value) {
			setSliceField(isPointer, schema, field, t)
		},
	}
}

func (defaults *Defaults) getTypeMappings(fieldType string) (TypeMappings, error) {
	if typeMappings, found := defaults.typeMappings[fieldType]; found {
		return typeMappings, nil
	} else {
		fmt.Println(typeMappings)
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
