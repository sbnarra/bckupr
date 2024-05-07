package filters

import (
	"slices"

	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/oapi/server"
)

func applyIncludeFilters(unfiltered map[string]*dockerTypes.Container, filters server.Filters) map[string]*dockerTypes.Container {
	if len(filters.IncludeNames) == 0 && len(filters.IncludeVolumes) == 0 {
		return unfiltered
	}

	filtered := make(map[string]*dockerTypes.Container)
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
