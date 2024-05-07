package api

import (
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

func (c *Client) Version() *errors.Error {
	return c.send("GET", PATH_VERSION, nil)
}

func (c *Client) Debug() *errors.Error {
	return c.send("GET", PATH_DEBUG, nil)
}

func (c *Client) CreateBackup(request *types.CreateBackupRequest) *errors.Error {
	return c.send("POST", PATH_BACKUPS, request)
}

func (c *Client) Rotate() *errors.Error {
	return c.send("POST", PATH_ROTATE, nil)
}

func (c *Client) List() *errors.Error {
	return c.send("GET", PATH_BACKUPS, nil)
}

func (c *Client) DeleteBackup(id string) *errors.Error {
	return c.send("DELETE", PATH_BACKUPS+"/"+id, nil)
}

func (c *Client) RestoreBackup(request *types.RestoreBackupRequest) *errors.Error {
	return c.send("POST", PATH_RESTORE_TRIGGER, request)
}

func (c *Client) BackupSchedule() *errors.Error {
	return c.send("GET", PATH_CRON_BACKUP_SCHEDULE, nil)
}
