package cobra

import (
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/spf13/cobra"
)

func RestoreBackupRequest(cmd *cobra.Command) (*types.RestoreBackupRequest, *errors.Error) {
	var err *errors.Error

	var taskArgs *types.TaskArgs
	if taskArgs, err = createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	}

	var notificationSettings *types.NotificationSettings
	if notificationSettings, err = createNotificationSettings(cmd); err != nil {
		return nil, err
	}

	return &types.RestoreBackupRequest{
		Args:                 *taskArgs,
		NotificationSettings: notificationSettings,
	}, nil
}

func createNotificationSettings(cmd *cobra.Command) (*types.NotificationSettings, *errors.Error) {
	var err *errors.Error

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

func DeleteRequest(cmd *cobra.Command) (*types.DeleteBackupRequest, *errors.Error) {
	if backupArgs, err := createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	} else {
		return &types.DeleteBackupRequest{
			Args: *backupArgs,
		}, nil
	}
}

func CreateBackupRequest(cmd *cobra.Command) (*types.CreateBackupRequest, *errors.Error) {
	var err *errors.Error

	var backupArgs *types.TaskArgs
	if backupArgs, err = createTaskArgs(keys.BackupStopModes, cmd); err != nil {
		return nil, err
	}

	var notificationSettings *types.NotificationSettings
	if notificationSettings, err = createNotificationSettings(cmd); err != nil {
		return nil, err
	}

	return &types.CreateBackupRequest{
		Args:                 *backupArgs,
		NotificationSettings: notificationSettings,
	}, nil
}

func createTaskArgs(stopModes *keys.Key, cmd *cobra.Command) (*types.TaskArgs, *errors.Error) {
	var err *errors.Error

	var filters *types.Filters
	if filters, err = createFilters(stopModes, cmd); err != nil {
		return nil, err
	}

	var labelPrefix string
	if labelPrefix, err = String(keys.LabelPrefix, cmd.Flags()); err != nil {
		return nil, err
	}

	var backupId string
	if backupId, err = String(keys.BackupId, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.TaskArgs{
		BackupId:    backupId,
		Filters:     *filters,
		LabelPrefix: labelPrefix,
	}, nil
}

func createFilters(stopModesKey *keys.Key, cmd *cobra.Command) (*types.Filters, *errors.Error) {
	var err *errors.Error

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

func DaemonInput(cmd *cobra.Command) (*types.DaemonInput, *errors.Error) {
	var err *errors.Error

	var backupDir string
	if backupDir, err = String(keys.HostBackupDir, cmd.Flags()); err != nil {
		return nil, err
	}

	var localContainersConfig string
	if localContainersConfig, err = String(keys.LocalContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	}

	var offsiteContainersConfig string
	if offsiteContainersConfig, err = String(keys.OffsiteContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	}

	var unixSocket string
	if unixSocket, err = String(keys.UnixSocket, cmd.Flags()); err != nil {
		return nil, err
	}

	var tcpAddr string
	if tcpAddr, err = String(keys.TcpAddr, cmd.Flags()); err != nil {
		return nil, err
	}

	var tcpApi bool
	if tcpApi, err = Bool(keys.TcpApi, cmd.Flags()); err != nil {
		return nil, err
	}

	var ui bool
	if ui, err = Bool(keys.UI, cmd.Flags()); err != nil {
		return nil, err
	}

	var metrics bool
	if metrics, err = Bool(keys.Metrics, cmd.Flags()); err != nil {
		return nil, err
	}

	var dockerHosts []string
	if dockerHosts, err = StringSlice(keys.DockerHosts, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.DaemonInput{
		BackupDir:               backupDir,
		DockerHosts:             dockerHosts,
		LocalContainersConfig:   localContainersConfig,
		OffsiteContainersConfig: offsiteContainersConfig,

		UnixSocket: unixSocket,
		TcpAddr:    tcpAddr,
		TcpApi:     tcpApi,
		UI:         ui,
		Metrics:    metrics,
	}, nil

}

func RotateBackupsRequest(cmd *cobra.Command) (*types.RotateBackupsRequest, *errors.Error) {
	var err *errors.Error

	var destroyBackups bool
	if destroyBackups, err = Bool(keys.DestroyBackups, cmd.Flags()); err != nil {
		return nil, err
	}

	var policyPath string
	if policyPath, err = String(keys.PoliciesPath, cmd.Flags()); err != nil {
		return nil, err
	}

	return &types.RotateBackupsRequest{
		Destroy:      destroyBackups,
		PoliciesPath: policyPath,
	}, nil
}
