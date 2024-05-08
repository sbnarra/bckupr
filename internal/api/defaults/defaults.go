package defaults

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func New(specPath string) (*Defaults, error) {
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromFile(specPath)
	return &Defaults{
		spec:         spec,
		typeMappings: basicTypeMappings,
		creators:     map[string]Creator{},
	}, err
}

func (d *Defaults) AddType(name string, new func() any) {
	d.creators[name] = Creator{
		New: new,
	}
}
