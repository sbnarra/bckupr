package client

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

type Client struct {
	client *spec.ClientWithResponses
}

func New(ctx contexts.Context, protocol string, server string) (*Client, *errors.Error) {
	url := protocol + "://" + server
	if http, err := spec.NewClientWithResponses(url, func(c *spec.Client) error {
		c.RequestEditors = []spec.RequestEditorFn{
			func(ct context.Context, req *http.Request) error {
				req.Header.Add("User-Agent", "bckupr-sdk/"+os.Getenv("VERSION"))
				req.Header.Add("Dry-Run", strconv.FormatBool(ctx.DryRun))
				req.Header.Add("Debug", strconv.FormatBool(ctx.Debug))
				return nil
			},
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create api client")
	} else {
		return &Client{http}, nil
	}
}

func (c *Client) TriggerBackup(ctx contexts.Context, req spec.ContainersConfig) (*spec.Backup, *errors.Error) {
	res, err := c.client.TriggerBackupWithResponse(ctx, req)
	if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error triggering backup")
	}
}

func (c *Client) TriggerBackupUsingId(ctx contexts.Context, id string, req spec.ContainersConfig) (*spec.Backup, *errors.Error) {
	res, err := c.client.TriggerBackupWithIdWithResponse(ctx, id, req)
	if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error triggering backup:"+id)
	}
}

func (c *Client) GetBackup(ctx contexts.Context, id string) (*spec.Backup, *errors.Error) {
	res, err := c.client.GetBackupWithResponse(ctx, id)
	if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error finding backup: "+id)
	}
}

func (c *Client) DeleteBackup(ctx contexts.Context, id string) *errors.Error {
	res, err := c.client.DeleteBackupWithResponse(ctx, id)
	if checkSuccess(res.HTTPResponse, err) {
		return nil
	} else {
		return errors.NewWrap(err, "error triggering backup")
	}
}

func (c *Client) ListBackups(ctx contexts.Context) ([]spec.Backup, *errors.Error) {
	res, err := c.client.ListBackupsWithResponse(ctx)
	if checkSuccess(res.HTTPResponse, err) {
		backups := res.JSON200
		return *backups, nil
	} else {
		return nil, errors.NewWrap(err, "error listing backups")
	}
}

func (c *Client) TriggerRestore(ctx contexts.Context, backupId string, req spec.ContainersConfig) (*spec.Restore, *errors.Error) {
	res, err := c.client.TriggerRestoreWithResponse(ctx, backupId, req)
	if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error triggering restore")
	}
}

func (c *Client) GetRestore(ctx contexts.Context, id string) (*spec.Restore, *errors.Error) {
	res, err := c.client.GetRestoreWithResponse(ctx, id)
	if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error finding restore: "+id)
	}
}

func (c *Client) RotateBackups(ctx contexts.Context, req spec.RotateInput) (*spec.Rotate, *errors.Error) {
	res, err := c.client.RotateBackupsWithResponse(ctx, req)
	if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error triggering rotate")
	}
}

func (c *Client) Version(ctx contexts.Context) (*spec.Version, *errors.Error) {
	res, err := c.client.GetVersionWithResponse(ctx)
	if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error getting version")
	}
}

func checkSuccess(res *http.Response, err error) bool {
	if err != nil {
		return false
	} else if (res.StatusCode / 100) != 2 {
		return false
	}
	return true
}
