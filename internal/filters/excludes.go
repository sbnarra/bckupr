package filters

import (
	"slices"

	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/oapi/server"
)

func applyExcludeFilters(unfiltered map[string]*dockerTypes.Container, filters server.Filters) map[string]*dockerTypes.Container {
	if len(filters.ExcludeNames) == 0 && len(filters.ExcludeVolumes) == 0 {
		return unfiltered
	}

	filtered := make(map[string]*dockerTypes.Container)
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
