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
		return nil, errors.New("nothing to " + name + " after applying include filters: names=" + strings.Join(filters.IncludeNames, ",") + ",volumes=" + strings.Join(filters.IncludeVolumes, ","))
	}

	filtered = applyExcludeFilters(filtered, filters)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + name + " after applying exclude filters: names=" + strings.Join(filters.ExcludeNames, ",") + ",volumes=" + strings.Join(filters.ExcludeVolumes, ","))
	}

	if stopModes != nil {
		filtered = applyStopModes(filtered, *stopModes)
	}
	if len(filtered) == 0 {
		stopModes := []string{}
		for _, stopMode := range stopModes {
			stopModes = append(stopModes, string(stopMode))
		}
		return nil, errors.New("nothing to " + name + " after applying stop modes: " + strings.Join(stopModes, ","))
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
