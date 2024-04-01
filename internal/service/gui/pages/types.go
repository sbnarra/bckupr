package pages
import (
	"time"

	"github.com/sbnarra/bckupr/pkg/types"
)

type IndexPage struct {
	Cron        Cron
	Backups     []types.Backup
	Error       error
	BackupInput *types.CreateBackupRequest
}

type Cron struct {
	Schedule string
	Next     time.Time
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
	Web       types.WebInput
}

type FeedbackPage struct {
	Action string
	Cron   Cron
	Error  error
}
