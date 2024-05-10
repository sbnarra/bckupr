package spec

import (
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

var (
	spec, _ = GetSwagger()
)

type D struct {
	entity string
}

func NewBackup() Backup {
	return Backup{}
}

var (
	BackupStopModes = []StopModes{
		Writers, Linked, Labelled,
	}
	RestoreStopModes = []StopModes{
		Attached, Linked, Labelled,
	}
)

func (c *ContainersConfig) WithDefaults(stopModes []StopModes) *errors.Error {
	d := D{entity: "ContainersConfig"}

	if c.LabelPrefix == nil || *c.LabelPrefix == "" {
		if labelPrefix, err := d.aString("label_prefix"); err != nil {
			return err
		} else {
			c.LabelPrefix = labelPrefix
		}
	}

	c.Filters.WithDefaults()
	return nil
}

func (r *RotateInput) WithDefaults() *errors.Error {
	return nil
}

func (f *Filters) WithDefaults() {
}

func (v *Version) NewVersion() *errors.Error {
	d := D{entity: "Version"}
	if v.Version == "" {
		if version, err := d.aString("version"); err != nil {
			return err
		} else {
			v.Version = *version
		}
	}
	return nil
}

func convert[T any](entity string, field string, empty T, conversion func(i any) T) (*T, *errors.Error) {
	if entitySchema, found := spec.Components.Schemas[entity]; !found {
		return nil, errors.Errorf("entity schema not found: entity=%v,field=%v", entity, field)
	} else if properties, found := entitySchema.Value.Properties[field]; !found {
		return nil, errors.Errorf("schmea properties not found: entity=%v,field=%v: %v", entity, field, encodings.ToJsonIE(entitySchema))
	} else {
		defaultV := properties.Value.Default
		if defaultV == nil {
			return &empty, nil
		}
		converted := conversion(defaultV)
		return &converted, nil
	}
}

func (d D) aString(field string) (*string, *errors.Error) {
	return convert(d.entity, field, "", func(i any) string {
		return i.(string)
	})
}

func (d D) aBool(field string) (*bool, *errors.Error) {
	return convert(d.entity, field, false, func(i any) bool {
		return i.(bool)
	})
}

func (d D) aInt(field string) (*int, *errors.Error) {
	return convert(d.entity, field, 0, func(i any) int {
		return i.(int)
	})
}

func (d D) aStringSlice(field string) (*[]string, *errors.Error) {
	return convert(d.entity, field, []string{}, func(i any) []string {
		return i.([]string)
	})
}
