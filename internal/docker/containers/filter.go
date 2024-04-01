package containers

import (
	"slices"

	"github.com/sbnarra/bckupr/pkg/types"
)

func ApplyFilters(unfiltered map[string]*Container, filters types.Filters) map[string]*Container {
	filtered := applyIncludeFilters(unfiltered, filters)
	filtered = applyExcludeFilters(filtered, filters)
	filtered = applyStopModes(filtered, filters.StopModes)
	return filtered
}

func applyIncludeFilters(unfiltered map[string]*Container, filters types.Filters) map[string]*Container {
	if len(filters.IncludeNames) == 0 && len(filters.IncludeVolumes) == 0 {
		return unfiltered
	}

	filtered := make(map[string]*Container)
	for id, container := range unfiltered {
		if len(filters.IncludeNames) != 0 {
			if slices.Contains(filters.IncludeNames, container.Name) {
				filtered[id] = container
			}
		}

		if len(filters.IncludeVolumes) != 0 {
			for name := range container.Volumes {
				if slices.Contains(filters.IncludeVolumes, name) {
					filtered[id] = container
				} else if backupsContain(filters.IncludeVolumes, container) {
					filtered[id] = container
				}
			}
		}
	}
	return filtered
}

func backupsContain(volumes []string, container *Container) bool {
	for name := range container.Backup.Volumes {
		if slices.Contains(volumes, name) {
			return true
		}
	}
	return false
}

func applyExcludeFilters(unfiltered map[string]*Container, filters types.Filters) map[string]*Container {
	if len(filters.ExcludeNames) == 0 && len(filters.ExcludeVolumes) == 0 {
		return unfiltered
	}

	filtered := make(map[string]*Container)
	for id, container := range filtered {

		if slices.Contains(filters.ExcludeNames, container.Name) {
			delete(filtered, id)
		}

		for name := range container.Volumes {
			if slices.Contains(filters.ExcludeVolumes, name) {
				delete(filtered, id)
			} else if !backupsContain(filters.ExcludeVolumes, container) {
				delete(filtered, id)
			}
		}
	}
	return filtered
}

func applyStopModes(unfiltered map[string]*Container, stopModes []string) map[string]*Container {
	if slices.Contains(stopModes, "all") {
		return unfiltered
	}

	backupPaths := []string{}
	for _, container := range unfiltered {
		for volume := range container.Backup.Volumes {
			backupPaths = append(backupPaths, volume)
		}
	}

	filtered := make(map[string]*Container)
	for id, container := range unfiltered {

		if slices.Contains(stopModes, "labelled") && container.Backup.Stop {
			filtered[id] = container
		}

		if slices.Contains(stopModes, "attached") {
			for _, path := range backupPaths {
				for volume := range container.Volumes {
					if path == volume {
						filtered[id] = container
					}
				}
			}
		}

		if slices.Contains(stopModes, "writers") {
			for _, path := range backupPaths {
				for volume, info := range container.Volumes {
					if info.Writer && path == volume {
						filtered[id] = container
					}
				}
			}
		}
	}
	return filtered
}
