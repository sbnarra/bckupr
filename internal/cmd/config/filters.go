package config

import (
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/spec"
	"github.com/spf13/cobra"
)

func initFilters(cmd *cobra.Command, stopModes *keys.Key) {
	flags.Register(keys.IncludeNames, cmd.Flags())
	flags.Register(keys.IncludeVolumes, cmd.Flags())
	flags.Register(keys.ExcludeName, cmd.Flags())
	flags.Register(keys.ExcludeVolumes, cmd.Flags())
}

func readFilters(cmd *cobra.Command) (*spec.Filters, *errors.Error) {
	var err *errors.Error

	var includeNames []string
	if includeNames, err = flags.StringSlice(keys.IncludeNames, cmd.Flags()); err != nil {
		return nil, err
	}

	var includeVolumes []string
	if includeVolumes, err = flags.StringSlice(keys.IncludeVolumes, cmd.Flags()); err != nil {
		return nil, err
	}

	var excludeNames []string
	if excludeNames, err = flags.StringSlice(keys.ExcludeName, cmd.Flags()); err != nil {
		return nil, err
	}

	var excludeVolumes []string
	if excludeVolumes, err = flags.StringSlice(keys.ExcludeVolumes, cmd.Flags()); err != nil {
		return nil, err
	}

	return &spec.Filters{
		ExcludeNames:   includeNames,
		ExcludeVolumes: includeVolumes,
		IncludeNames:   excludeNames,
		IncludeVolumes: excludeVolumes,
	}, nil
}
