package server

var (
	spec, _ = GetSwagger()
)

type D struct {
	entity string
}

func NewTriggerBackup() TriggerBackup {
	trigger := TriggerBackup{}
	trigger.FromTask(NewTask([]StopModes{
		Writers, Linked, Labelled,
	}))
	return trigger
}

func NewTriggerRestore() TriggerRestore {
	trigger := TriggerRestore{}
	trigger.FromTask(NewTask([]StopModes{
		Attached, Linked, Labelled,
	}))
	return trigger
}

func NewTask(stopModes []StopModes) Task {
	d := D{entity: "Task"}
	return Task{
		Filters:     NewFilters(),
		LabelPrefix: *d.aString("LabelPrefix"),
		StopModes:   stopModes,
	}
}

func NewFilters() Filters {
	return Filters{
		IncludeNames:   []string{},
		IncludeVolumes: []string{},
		ExcludeNames:   []string{},
		ExcludeVolumes: []string{},
	}
}

func NewVersion() Version {
	d := D{entity: "Version"}
	return Version{
		Version: *d.aString("Version"),
	}
}

func convert[T any](entity string, field string, empty T, conversion func(i any) T) *T {
	origin := spec.Components.Schemas[entity].Value.Properties[field].Value.Default
	if origin == nil {
		return &empty
	}
	converted := conversion(origin)
	return &converted
}

func (d D) aString(field string) *string {
	return convert(d.entity, field, "", func(i any) string {
		return i.(string)
	})
}

func (d D) aBool(field string) *bool {
	return convert(d.entity, field, false, func(i any) bool {
		return i.(bool)
	})
}

func (d D) aInt(field string) *int {
	return convert(d.entity, field, 0, func(i any) int {
		return i.(int)
	})
}

func (d D) aStringSlice(field string) *[]string {
	return convert(d.entity, field, []string{}, func(i any) []string {
		return i.([]string)
	})
}
