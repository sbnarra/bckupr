package config

import (
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/spec"
	"github.com/spf13/cobra"
)

func InitTaskTrigger(cmd *cobra.Command, stopModes *keys.Key) {
	initFilters(cmd)

	flags.Register(stopModes, cmd.Flags())
	flags.Register(keys.LabelPrefix, cmd.Flags())
	flags.Register(keys.BackupId, cmd.Flags())
}

func ReadContainersConfig(cmd *cobra.Command, stopModesKey *keys.Key) (string, *spec.ContainersConfig, *errors.Error) {
	stopModes := []spec.StopModes{}
	if stopModesS, err := flags.StringSlice(stopModesKey, cmd.Flags()); err != nil {
		return "", nil, err
	} else {
		for _, stopMode := range stopModesS {
			stopModes = append(stopModes, spec.StopModes(stopMode))
		}
	}

	if filters, err := readFilters(cmd); err != nil {
		return "", nil, err
	} else if labelPrefix, err := flags.String(keys.LabelPrefix, cmd.Flags()); err != nil {
		return "", nil, err
	} else if backupId, err := flags.String(keys.BackupId, cmd.Flags()); err != nil {
		return "", nil, err
	} else {
		return backupId, &spec.ContainersConfig{
			Filters:     *filters,
			LabelPrefix: &labelPrefix,
			StopModes:   &stopModes,
		}, nil
	}
}
