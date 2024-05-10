package defaults

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func Fill(i any, defaults *Defaults) error {
	return fillUsingName(fmt.Sprintf("%T", i), i, defaults)
}

func fillUsingName(schemaName string, i any, defaults *Defaults) error {
	fmt.Println("schema name is", schemaName)
	if strings.Contains(schemaName, ".") {
		schemaName = strings.Split(schemaName, ".")[1]
		fmt.Println("schema name is now", schemaName)
	}

	schema, found := defaults.spec.Components.Schemas[schemaName]
	if !found {
		return fmt.Errorf(
			"add missing creator: Defaults.AddType(\"%v\", func() interface{} { return &%v{} })",
			schemaName, schemaName)
	}

	vo := reflect.ValueOf(i)
	el := vo.Elem()

	for name, ref := range schema.Value.Properties {
		fmt.Println("Setting Field", strings.Title(name))
		field := el.FieldByName(strings.Title(name))
		if err := setValue(ref.Value, field, defaults); err != nil {
			return err
		}
	}
	return nil
}

func setValue(schema *openapi3.Schema, field reflect.Value, defaults *Defaults) error {
	k := field.Kind()
	if k == reflect.Pointer {
		k = field.Elem().Kind()
	}
	fmt.Println("setValue", field.Type(), k)

	switch field.Kind() {
	case reflect.Struct:
		if creator, found := defaults.creators[field.Type().String()]; !found {
			return fmt.Errorf("add missing creator: Defaults.AddType(\"%v\", func() any { return &%v{} })",
				field.Type().String(),
				field.Type().String())
		} else {
			instance := creator.New()
			el := reflect.ValueOf(instance)
			el = el.Elem()
			fillUsingName(field.Type().String(), instance, defaults)
			field.Set(el)
			return nil
		}
	}

	if schema.Default == nil {
		fmt.Println("no default", field.Type())
		return nil
	}

	fieldType := field.Type().String()
	isPointer := false
	isSlice := false

	if strings.HasPrefix(fieldType, "*") {
		fieldType = fieldType[1:]
		isPointer = true
	}

	if strings.HasPrefix(fieldType, "[]") {
		fieldType = fieldType[2:]
		isSlice = true
	}

	var err error
	var typeMappings TypeMappings
	if typeMappings, err = defaults.getTypeMappings(fieldType); err != nil {
		return err
	}

	if isSlice {
		typeMappings.Slice(isPointer, schema, field)
	} else {
		fmt.Println("before:", field)
		typeMappings.Single(isPointer, schema, field)
		fmt.Println("after:", field)
	}

	return nil
}
