package pages

import (
	"github.com/sbnarra/bckupr/pkg/types"
)

type IndexPage struct {
	Cron        Cron
	Backups     []*types.Backup
	Error       error
	BackupInput *types.CreateBackupRequest
}

type Cron struct {
	BackupSchedule string
	NextBackup     string
	RotateSchedule string
	NextRotate     string
}

type SettingsPage struct {
	Cron          Cron
	Global        GlobalSettings
	Notifications *types.NotificationSettings
	Error         error
}

type GlobalSettings struct {
	DryRun    bool
	Debug     bool
	BackupDir string
	Args      types.TaskArgs
	Web       types.DaemonInput
}

type FeedbackPage struct {
	Action string
	Cron   Cron
	Error  error
}
