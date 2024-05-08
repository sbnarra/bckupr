package filters

import (
	"slices"
	"strings"

	"github.com/sbnarra/bckupr/internal/api/spec"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func Apply(ctx contexts.Context, unfiltered map[string]*dockerTypes.Container, filters spec.Filters, stopModes []spec.StopModes) (map[string]*dockerTypes.Container, *errors.Error) {
	filtered := applyIncludeFilters(unfiltered, filters)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + ctx.Name + " after applying include filters: names=" + strings.Join(filters.IncludeNames, ",") + ",volumes=" + strings.Join(filters.IncludeVolumes, ","))
	}

	filtered = applyExcludeFilters(filtered, filters)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + ctx.Name + " after applying exclude filters: names=" + strings.Join(filters.ExcludeNames, ",") + ",volumes=" + strings.Join(filters.ExcludeVolumes, ","))
	}

	filtered = applyStopModes(filtered, stopModes)
	if len(filtered) == 0 {
		stopModes := []string{}
		for _, stopMode := range stopModes {
			stopModes = append(stopModes, string(stopMode))
		}
		return nil, errors.New("nothing to " + ctx.Name + " after applying stop modes: " + strings.Join(stopModes, ","))
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
