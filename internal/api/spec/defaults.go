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

func (b *BackupTrigger) WithDefaults() *errors.Error {
	if task, err := b.AsTaskTrigger(); err != nil {
		return errors.Wrap(err, "error unwrapping task")
	} else {
		task.WithDefaults([]StopModes{
			Writers, Linked, Labelled,
		})
		b.FromTaskTrigger(task)
	}
	return nil
}

func (r *RestoreTrigger) WithDefaults() *errors.Error {
	if task, err := r.AsTaskTrigger(); err != nil {
		return errors.Wrap(err, "error unwrapping task")
	} else {
		task.WithDefaults([]StopModes{
			Attached, Linked, Labelled,
		})
		r.FromTaskTrigger(task)
	}
	return nil
}

func (r *RotateTrigger) WithDefaults() *errors.Error {
	return nil
}

func (t *TaskTrigger) WithDefaults(stopModes []StopModes) *errors.Error {
	d := D{entity: "TaskTrigger"}

	if t.LabelPrefix == nil || *t.LabelPrefix == "" {
		if labelPrefix, err := d.aString("label_prefix"); err != nil {
			return err
		} else {
			t.LabelPrefix = labelPrefix
		}
	}

	t.Filters.WithDefaults()
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
