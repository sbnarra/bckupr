package filters

import (
	"slices"
	"strings"

	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func Apply(ctx contexts.Context, unfiltered map[string]*dockerTypes.Container, filters publicTypes.Filters) (map[string]*dockerTypes.Container, *errors.Error) {
	filtered := applyIncludeFilters(unfiltered, filters)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + ctx.Name + " after applying include filters: names=" + strings.Join(filters.IncludeNames, ",") + ",volumes=" + strings.Join(filters.IncludeVolumes, ","))
	}

	filtered = applyExcludeFilters(filtered, filters)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + ctx.Name + " after applying exclude filters: names=" + strings.Join(filters.ExcludeNames, ",") + ",volumes=" + strings.Join(filters.ExcludeVolumes, ","))
	}

	filtered = applyStopModes(filtered, filters.StopModes)
	if len(filtered) == 0 {
		return nil, errors.New("nothing to " + ctx.Name + " after applying stop modes: " + strings.Join(filters.StopModes, ","))
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
