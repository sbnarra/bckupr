package filters

import (
	"context"
	"slices"
	"strings"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/contexts"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func Apply(ctx context.Context, unfiltered map[string]*dockerTypes.Container, filters spec.Filters, stopModes *[]spec.StopModes) (map[string]*dockerTypes.Container, *errors.E) {
	name := contexts.Name(ctx)
	filtered := applyIncludeFilters(unfiltered, filters)
	if len(filtered) == 0 {
		return nil, errors.Errorf("nothing to %v after applying include filters: names=%v,volumes=%v", name, strings.Join(filters.IncludeNames, ","), strings.Join(filters.IncludeVolumes, ","))
	}

	filtered = applyExcludeFilters(filtered, filters)
	if len(filtered) == 0 {
		return nil, errors.Errorf("nothing to %v after applying exclude filters: names=%v,volumes=%v", name, strings.Join(filters.ExcludeNames, ","), strings.Join(filters.ExcludeVolumes, ","))
	}

	if stopModes != nil {
		filtered = applyStopModes(filtered, *stopModes)
	}
	if len(filtered) == 0 {
		stopModes := []string{}
		for _, stopMode := range stopModes {
			stopModes = append(stopModes, string(stopMode))
		}
		return nil, errors.Errorf("nothing to %v after applying stop modes: %v", name, strings.Join(stopModes, ","))
	}

	return filtered, nil
}

func backupsContain(volumes []string, container *dockerTypes.Container) bool {
	for name := range container.Backup.Volumes {
		if slices.Contains(volumes, name) {
			return true
		}
	}
	return false
}
