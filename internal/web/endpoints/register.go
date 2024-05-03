package endpoints

import (
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/web/dispatcher"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func Register(d *dispatcher.Dispatcher, cron *cron.Cron, socket string, containers types.ContainerTemplates) {
	d.POST(api.PATH_BACKUPS, createBackup(containers))
	d.GET(api.PATH_BACKUPS, listBackups)
	d.DELETE(api.PATH_BACKUPS, deleteBackup)

	d.POST(api.PATH_RESTORE_TRIGGER, restoreBackup(containers))

	d.GET(api.PATH_CRON_BACKUP_SCHEDULE, backupSchedule(cron))

	d.GET(api.PATH_DEBUG, debug(socket))
	d.GET(api.PATH_VERSION, version)

}
