package list

import (
	"slices"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/sbnarra/bckupr/internal/docker/types"
)

func createBackupConfig(container dockerTypes.Container, labelPrefix string) types.BackupConfig {
	labelKeys := make([]string, 0)
	volumes := make(map[string]string)

	for key, value := range container.Labels {

		if strings.HasPrefix(key, labelPrefix+".") {
			labelKeys = append(labelKeys, key)
		} else {
			continue
		}

		if labelPrefix+".volumes" == key {
			for _, name := range strings.Split(value, ",") {
				volumes[name] = name
			}
		} else if strings.HasPrefix(key, labelPrefix+".volumes.") {
			key = key[len(labelPrefix+".volumes."):]
			volumes[key] = value
		}
	}

	return types.BackupConfig{
		Ignore:     isLabelTrue(container.Labels, labelPrefix, labelKeys, "ignore"),
		Stop:       isLabelTrue(container.Labels, labelPrefix, labelKeys, "stop"),
		Filesystem: isLabelTrue(container.Labels, labelPrefix, labelKeys, "filesystem"),
		Volumes:    volumes,
	}
}

func isLabelTrue(labels map[string]string, labelPrefix string, labelKeys []string, key string) bool {
	return slices.Contains(labelKeys, labelPrefix+"."+key) &&
		strings.ToLower(labels[labelPrefix+"."+key]) == "true"
}
