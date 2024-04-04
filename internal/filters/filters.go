package filters

import (
	"slices"

	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func Apply(unfiltered map[string]*dockerTypes.Container, filters publicTypes.Filters) map[string]*dockerTypes.Container {
	filtered := applyIncludeFilters(unfiltered, filters)
	filtered = applyExcludeFilters(filtered, filters)
	filtered = applyStopModes(filtered, filters.StopModes)
	return filtered
}

func backupsContain(volumes []string, container *dockerTypes.Container) bool {
	for name := range container.Backup.Volumes {
		if slices.Contains(volumes, name) {
			return true
		}
	}
	return false
}
