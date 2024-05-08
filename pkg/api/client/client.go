package client

import (
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

type Client struct {
	client *spec.ClientWithResponses
}

func New(server string) (*Client, *errors.Error) {
	if http, err := spec.NewClientWithResponses(server); err != nil {
		return nil, errors.Wrap(err, "failed to create api client")
	} else {
		return &Client{http}, nil
	}
}

func (c *Client) TriggerBackup(ctx contexts.Context, req spec.BackupTrigger) (*spec.Task, *spec.Backup, *errors.Error) {
	return c.TriggerBackupUsingId(ctx, "", req)
}

func (c *Client) TriggerBackupUsingId(ctx contexts.Context, id string, req spec.BackupTrigger) (*spec.Task, *spec.Backup, *errors.Error) {
	res, err := c.client.TriggerBackupWithResponse(ctx, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error triggering backup")
	}
	backup := res.JSON200
	task, err := backup.AsTask()
	return &task, backup, errors.Wrap(err, "error reading task")
}

func (c *Client) TriggerRestore(ctx contexts.Context, backupId string, req spec.RestoreTrigger) (*spec.Task, *errors.Error) {
	res, err := c.client.TriggerRestoreWithResponse(ctx, backupId, req)
	if err != nil {
		return nil, errors.Wrap(err, "error triggering restore")
	}
	restore := res.JSON200
	return restore, nil
}

func (c *Client) DeleteBackup(ctx contexts.Context, id string) *errors.Error {
	_, err := c.client.DeleteBackupWithResponse(ctx, id)
	return errors.Wrap(err, "error triggering backup")
}

func (c *Client) ListBackups(ctx contexts.Context) ([]spec.Backup, *errors.Error) {
	res, err := c.client.ListBackupsWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error listing backups")
	}
	backups := res.JSON200
	return *backups, nil
}

func (c *Client) RotateBackups(ctx contexts.Context, req spec.RotateTrigger) (*spec.Task, *errors.Error) {
	res, err := c.client.RotateBackupsWithResponse(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "error triggering rotate")
	}
	task := res.JSON200
	return task, nil
}

func (c *Client) Version(ctx contexts.Context) (*spec.Version, *errors.Error) {
	res, err := c.client.GetVersionWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting version")
	}
	task := res.JSON200
	return task, nil
}
