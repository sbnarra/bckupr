package defaults

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
)

func New() (*Defaults, error) {
	spec, err := spec.GetSwagger()

	return &Defaults{
		spec:         spec,
		typeMappings: basicTypeMappings,
		creators:     map[string]Creator{},
	}, err
}

func (d *Defaults) AddCreator(name string, new func() any) {
	d.creators[name] = Creator{
		New: new,
	}
}

func (d *Defaults) AddTypeMapping(name string, mapping func(any) any) {
	d.typeMappings[name] = standardTypeMapping(mapping)
}
