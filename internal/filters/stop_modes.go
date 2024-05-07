package filters

import (
	"slices"

	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/oapi/server"
)

func applyStopModes(unfiltered map[string]*dockerTypes.Container, stopModes []server.StopModes) map[string]*dockerTypes.Container {
	if slices.Contains(stopModes, server.All) {
		return unfiltered
	}

	backupPaths := []string{}
	for _, container := range unfiltered {
		for volume := range container.Backup.Volumes {
			backupPaths = append(backupPaths, volume)
		}
	}

	filtered := make(map[string]*dockerTypes.Container)
	for id, container := range unfiltered {

		if slices.Contains(stopModes, server.Labelled) && container.Backup.Stop {
			filtered[id] = container
		}

		if slices.Contains(stopModes, server.Attached) {
			for _, path := range backupPaths {
				for volume := range container.Volumes {
					if path == volume {
						filtered[id] = container
					}
				}
			}
		}

		if slices.Contains(stopModes, server.Writers) {
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
