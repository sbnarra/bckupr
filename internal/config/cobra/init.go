package cobra

import (
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/spf13/cobra"
)

func InitGlobal(cmd *cobra.Command) {
	register(keys.DryRun, cmd.PersistentFlags())
	register(keys.Debug, cmd.PersistentFlags())
}

func InitDaemon(cmd *cobra.Command) {
	InitCron(cmd)

	register(keys.UnixSocket, cmd.Flags())
	register(keys.TcpAddr, cmd.Flags())
	register(keys.ExposeApi, cmd.Flags())
	register(keys.UiEnabled, cmd.Flags())
	register(keys.MetricsEnabled, cmd.Flags())
}

func InitCron(cmd *cobra.Command) {
	InitBackup(cmd)

	register(keys.BackupSchedule, cmd.Flags())
	register(keys.TimeZone, cmd.Flags())
}

func InitDebug(cmd *cobra.Command) {
	InitDaemonClient(cmd)
}

func InitList(cmd *cobra.Command) {
	InitDaemonClient(cmd)
	initTaskArgs(cmd, keys.BackupStopModes)
}

func InitDelete(cmd *cobra.Command) {
	InitDaemonClient(cmd)
	initTaskArgs(cmd, keys.BackupStopModes)

	register(keys.BackupId, cmd.Flags())
}

func InitBackup(cmd *cobra.Command) {
	InitDaemonClient(cmd)
	initTaskArgs(cmd, keys.BackupStopModes)

	register(keys.BackupIdOverride, cmd.Flags())
}

func InitDaemonClient(cmd *cobra.Command) {
	register(keys.NoDaemon, cmd.Flags())
	register(keys.DaemonAddr, cmd.Flags())
	register(keys.DaemonProtocol, cmd.Flags())
	register(keys.DaemonNet, cmd.Flags())
}

func InitRestore(cmd *cobra.Command) {
	InitDaemonClient(cmd)

	register(keys.BackupId, cmd.Flags())
	required(keys.BackupId, cmd)

	initTaskArgs(cmd, keys.BackupStopModes)
}

func initTaskArgs(cmd *cobra.Command, stopModes *keys.Key) {
	initFilters(cmd, stopModes)
	initNotifications(cmd)

	register(keys.DockerHosts, cmd.Flags())
	register(keys.LabelPrefix, cmd.Flags())

	register(keys.BackupDir, cmd.Flags())
	required(keys.BackupDir, cmd)

	register(keys.LocalContainers, cmd.Flags())
	register(keys.OffsiteContainers, cmd.Flags())
}

func initNotifications(cmd *cobra.Command) {
	register(keys.NotificationUrls, cmd.Flags())

	register(keys.NotifyJobStarted, cmd.Flags())
	register(keys.NotifyJobCompleted, cmd.Flags())
	register(keys.NotifyJobError, cmd.Flags())

	register(keys.NotifyTaskStarted, cmd.Flags())
	register(keys.NotifyTaskCompleted, cmd.Flags())
	register(keys.NotifyTaskError, cmd.Flags())
}

func initFilters(cmd *cobra.Command, stopModes *keys.Key) {
	register(stopModes, cmd.Flags())
	register(keys.IncludeNames, cmd.Flags())
	register(keys.IncludeVolumes, cmd.Flags())
	register(keys.ExcludeName, cmd.Flags())
	register(keys.ExcludeVolumes, cmd.Flags())
}
