package tasks

import (
	"slices"

	"github.com/sbnarra/bckupr/internal/api/spec"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
)

func convertToTasks(containerList map[string]*dockerTypes.Container, filters spec.Filters) map[string]*task {
	tasks := make(map[string]*task)
	for _, container := range containerList {
		for name, path := range container.Backup.Volumes {

			if addTask(container.Name, name, filters) {
				tasks[name] = &task{
					Completed:  false,
					Volume:     path,
					Containers: []*dockerTypes.Container{},
				}
			}
		}
	}

	populateTaskContainers(tasks, containerList)
	return tasks
}

func addTask(conName string, volName string, filters spec.Filters) bool {
	if len(filters.IncludeNames) != 0 || len(filters.IncludeVolumes) != 0 || len(filters.ExcludeNames) != 0 || len(filters.ExcludeVolumes) != 0 {

		in := len(filters.IncludeNames) != 0 && slices.Contains(filters.IncludeNames, conName)
		iv := len(filters.IncludeVolumes) != 0 && slices.Contains(filters.IncludeVolumes, volName)

		en := len(filters.ExcludeNames) != 0 && slices.Contains(filters.ExcludeNames, conName)
		ev := len(filters.ExcludeVolumes) != 0 && slices.Contains(filters.ExcludeVolumes, volName)

		if en || ev {
			return false
		}
		return in || iv
	}
	return true
}

func populateTaskContainers(tasks map[string]*task, containerList map[string]*dockerTypes.Container) {
	for _, task := range tasks {

		containersMatchingVolume := make(map[string]*dockerTypes.Container)
		for id, container := range containerList {
			for volume := range container.Volumes {
				// windows/macOS uses a VM which mounts OS directories to the VM under /host_mnt
				// clearly this is going to break for someone, guessing people could change
				if volume == task.Volume || volume == "/host_mnt"+task.Volume {
					containersMatchingVolume[id] = container
				}
			}
		}

		for _, container := range containersMatchingVolume {
			task.Containers = append(task.Containers, container)
		}
	}
}
