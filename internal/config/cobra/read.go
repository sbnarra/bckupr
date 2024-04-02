package cobra

import (
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/spf13/cobra"
)

func RestoreBackupRequest(cmd *cobra.Command) (*types.RestoreBackupRequest, error) {
	var err error

	var taskArgs *types.TaskArgs
	if taskArgs, err = createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	}

	var backupId string
	if backupId, err = String(keys.BackupId, cmd.Flags()); err != nil {
		return nil, err
	}

	var notificationSettings *types.NotificationSettings
	if notificationSettings, err = createNotificationSettings(cmd); err != nil {
		return nil, err
	}

	return &types.RestoreBackupRequest{
		BackupId:             backupId,
		Args:                 *taskArgs,
		NotificationSettings: notificationSettings,
	}, nil
}

func createNotificationSettings(cmd *cobra.Command) (*types.NotificationSettings, error) {
	var err error

	var notificationUrls []string
	if notificationUrls, err = StringSlice(keys.NotificationUrls, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyJobStarted bool
	if notifyJobStarted, err = Bool(keys.NotifyJobStarted, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyJobCompleted bool
	if notifyJobCompleted, err = Bool(keys.NotifyJobCompleted, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyJobError bool
	if notifyJobError, err = Bool(keys.NotifyJobError, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyTaskStarted bool
	if notifyTaskStarted, err = Bool(keys.NotifyTaskStarted, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyTaskCompleted bool
	if notifyTaskCompleted, err = Bool(keys.NotifyTaskCompleted, cmd.Flags()); err != nil {
		return nil, err
	}

	var notifyTaskError bool
	if notifyTaskError, err = Bool(keys.NotifyTaskError, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.NotificationSettings{
		NotificationUrls: notificationUrls,

		NotifyJobStarted:    notifyJobStarted,
		NotifyJobCompleted:  notifyJobCompleted,
		NotifyJobError:      notifyJobError,
		NotifyTaskStarted:   notifyTaskStarted,
		NotifyTaskCompleted: notifyTaskCompleted,
		NotifyTaskError:     notifyTaskError,
	}, nil
}

func DeleteRequest(cmd *cobra.Command) (*types.DeleteBackupRequest, error) {
	if backupArgs, err := createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	} else if backupId, err := String(keys.BackupId, cmd.Flags()); err != nil {
		return nil, err
	} else {
		return &types.DeleteBackupRequest{
			BackupId: backupId,
			Args:     *backupArgs,
		}, nil
	}
}

func CreateBackupRequest(cmd *cobra.Command) (*types.CreateBackupRequest, error) {
	var err error

	var backupArgs *types.TaskArgs
	if backupArgs, err = createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	}

	var backupIdOverride string
	if backupIdOverride, err = String(keys.BackupIdOverride, cmd.Flags()); err != nil {
		return nil, err
	}

	var notificationSettings *types.NotificationSettings
	if notificationSettings, err = createNotificationSettings(cmd); err != nil {
		return nil, err
	}

	return &types.CreateBackupRequest{
		BackupIdOverride:     backupIdOverride,
		Args:                 *backupArgs,
		NotificationSettings: notificationSettings,
	}, nil
}

func createTemplate(name string, cmd *cobra.Command) (*types.ContainerTemplate, error) {
	var err error

	var image string
	if image, err = cmd.Flags().GetString(name + "-image"); err != nil {
		return nil, err
	}

	var command []string
	if command, err = cmd.Flags().GetStringSlice(name + "-cmd"); err != nil {
		return nil, err
	}

	var env []string
	if env, err = cmd.Flags().GetStringSlice(name + "-env"); err != nil {
		return nil, err
	}

	return &types.ContainerTemplate{
		Alias:   name,
		Image:   image,
		Cmd:     command,
		Env:     env,
		Volumes: []string{},
	}, nil
}

func createTaskArgs(stopModes *keys.Key, cmd *cobra.Command) (*types.TaskArgs, error) {
	var err error

	var dockerHosts []string
	if dockerHosts, err = StringSlice(keys.DockerHosts, cmd.Flags()); err != nil {
		return nil, err
	}

	var filters *types.Filters
	if filters, err = createFilters(stopModes, cmd); err != nil {
		return nil, err
	}

	var labelPrefix string
	if labelPrefix, err = String(keys.LabelPrefix, cmd.Flags()); err != nil {
		return nil, err
	}

	var localContainersConfig string
	if localContainersConfig, err = String(keys.LocalContainers, cmd.Flags()); err != nil {
		return nil, err
	}

	var offsiteContainersConfig string
	if offsiteContainersConfig, err = String(keys.OffsiteContainers, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.TaskArgs{
		DockerHosts:             dockerHosts,
		Filters:                 *filters,
		LabelPrefix:             labelPrefix,
		LocalContainersConfig:   localContainersConfig,
		OffsiteContainersConfig: offsiteContainersConfig,
	}, nil
}

func createFilters(stopModesKey *keys.Key, cmd *cobra.Command) (*types.Filters, error) {
	var err error

	var stopModes []string
	if stopModes, err = StringSlice(stopModesKey, cmd.Flags()); err != nil {
		return nil, err
	}

	var includeNames []string
	if includeNames, err = StringSlice(keys.IncludeNames, cmd.Flags()); err != nil {
		return nil, err
	}

	var includeVolumes []string
	if includeVolumes, err = StringSlice(keys.IncludeVolumes, cmd.Flags()); err != nil {
		return nil, err
	}

	var excludeNames []string
	if excludeNames, err = StringSlice(keys.ExcludeName, cmd.Flags()); err != nil {
		return nil, err
	}

	var excludeVolumes []string
	if excludeVolumes, err = StringSlice(keys.ExcludeVolumes, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.Filters{
		StopModes:      stopModes,
		IncludeNames:   includeNames,
		IncludeVolumes: includeVolumes,
		ExcludeNames:   excludeNames,
		ExcludeVolumes: excludeVolumes,
	}, nil
}
