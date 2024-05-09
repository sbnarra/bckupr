package defaults

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
)

func New(specPath string) (*Defaults, error) {
	spec, err := spec.GetSwagger()

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
