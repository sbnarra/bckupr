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

func (c *TaskInput) WithDefaults(stopModes []StopModes) *errors.E {
	d := D{entity: "TaskInput"}

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

func (t *TaskInput) IsDryRun() bool {
	if t == nil || t.NoDryRun == nil {
		return true
	}
	noDryRun := *t.NoDryRun
	return !noDryRun
}

func (r *RotateInput) IsDryRun() bool {
	if r == nil || r.NoDryRun == nil {
		return true
	}
	noDryRun := *r.NoDryRun
	return !noDryRun
}

func (r *RotateInput) WithDefaults() *errors.E {
	return nil
}

func (f *Filters) WithDefaults() {
}

func (v *Version) NewVersion() *errors.E {
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

func convert[T any](entity string, field string, empty T, conversion func(i any) T) (*T, *errors.E) {
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

func (d D) aString(field string) (*string, *errors.E) {
	return convert(d.entity, field, "", func(i any) string {
		return i.(string)
	})
}

func (d D) aBool(field string) (*bool, *errors.E) {
	return convert(d.entity, field, false, func(i any) bool {
		return i.(bool)
	})
}

func (d D) aInt(field string) (*int, *errors.E) {
	return convert(d.entity, field, 0, func(i any) int {
		return i.(int)
	})
}

func (d D) aStringSlice(field string) (*[]string, *errors.E) {
	return convert(d.entity, field, []string{}, func(i any) []string {
		return i.([]string)
	})
}
