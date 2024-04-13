package api

import "github.com/sbnarra/bckupr/pkg/types"

func (c *Client) Version() error {
	return c.send("GET", PATH_VERSION, nil)
}

func (c *Client) Debug() error {
	return c.send("GET", PATH_DEBUG, nil)
}

func (c *Client) Backup(request *types.CreateBackupRequest) error {
	return c.send("POST", PATH_BACKUPS, request)
}

func (c *Client) Rotate() error {
	return c.send("POST", PATH_ROTATE, nil)
}

func (c *Client) List() error {
	return c.send("GET", PATH_BACKUPS, nil)
}

func (c *Client) Delete(request *types.DeleteBackupRequest) error {
	return c.send("DELETE", PATH_BACKUPS, request)
}

func (c *Client) Restore(request *types.RestoreBackupRequest) error {
	return c.send("POST", PATH_RESTORE_TRIGGER, request)
}

func (c *Client) BackupSchedule() error {
	return c.send("GET", PATH_CRON_BACKUP_SCHEDULE, nil)
}
