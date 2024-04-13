package types

import "github.com/sbnarra/bckupr/internal/config/keys"

func DefaultCreateBackupRequest() *CreateBackupRequest {
	return &CreateBackupRequest{
		DryRun:               keys.DryRun.EnvBool(),
		Args:                 DefaultTaskArgs(keys.BackupStopModes),
		NotificationSettings: DefaultNotificationSettings(),
	}
}

func DefaultDeleteBackupRequest() *DeleteBackupRequest {
	return &DeleteBackupRequest{
		Args: DefaultTaskArgs(keys.BackupStopModes),
	}
}

func DefaultRestoreBackupRequest() *RestoreBackupRequest {
	return &RestoreBackupRequest{
		DryRun:               keys.DryRun.EnvBool(),
		Args:                 DefaultTaskArgs(keys.RestoreStopModes),
		NotificationSettings: DefaultNotificationSettings(),
	}
}

func DefaultNotificationSettings() *NotificationSettings {
	return &NotificationSettings{
		NotificationUrls:    keys.NotificationUrls.EnvStringSlice(),
		NotifyJobStarted:    keys.NotifyJobStarted.EnvBool(),
		NotifyJobCompleted:  keys.NotifyJobCompleted.EnvBool(),
		NotifyJobError:      keys.NotifyJobError.EnvBool(),
		NotifyTaskStarted:   keys.NotifyTaskStarted.EnvBool(),
		NotifyTaskCompleted: keys.NotifyTaskCompleted.EnvBool(),
		NotifyTaskError:     keys.NotifyTaskError.EnvBool(),
	}
}

func DefaultTaskArgs(stopModes *keys.Key) TaskArgs {
	return TaskArgs{
		BackupId:                keys.BackupId.EnvString(),
		DockerHosts:             keys.DockerHosts.EnvStringSlice(),
		Filters:                 defaultFilters(stopModes),
		LabelPrefix:             keys.LabelPrefix.EnvString(),
		LocalContainersConfig:   keys.LocalContainers.EnvString(),
		OffsiteContainersConfig: keys.OffsiteContainers.EnvString(),
	}
}

func defaultFilters(stopModes *keys.Key) Filters {
	return Filters{
		StopModes:      stopModes.EnvStringSlice(),
		IncludeNames:   []string{},
		IncludeVolumes: []string{},
		ExcludeNames:   []string{},
		ExcludeVolumes: []string{},
	}
}

func DefaultWebInput() WebInput {
	return WebInput{
		UnixSocket:     keys.UnixSocket.EnvString(),
		TcpAddr:        keys.TcpAddr.EnvString(),
		ExposeApi:      keys.ExposeApi.EnvBool(),
		UiEnabled:      keys.UiEnabled.EnvBool(),
		MetricsEnabled: keys.MetricsEnabled.EnvBool(),
	}
}
