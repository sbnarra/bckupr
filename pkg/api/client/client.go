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

func (c *Client) StartBackup(ctx contexts.Context, req spec.ContainersConfig) (*spec.Backup, *errors.Error) {
	res, err := c.client.StartBackupWithResponse(ctx, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error starting backup: "+string(res.Body))
	}
}

func (c *Client) StartBackupWithId(ctx contexts.Context, id string, req spec.ContainersConfig) (*spec.Backup, *errors.Error) {
	res, err := c.client.StartBackupWithIdWithResponse(ctx, id, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error starting backup:"+id+": "+string(res.Body))
	}
}

func (c *Client) GetBackup(ctx contexts.Context, id string) (*spec.Backup, *errors.Error) {
	res, err := c.client.GetBackupWithResponse(ctx, id)
	if res == nil {
		return nil, errors.NewWrap(err, "error finding backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error finding backup: "+id+": "+string(res.Body))
	}
}

func (c *Client) DeleteBackup(ctx contexts.Context, id string) *errors.Error {
	res, err := c.client.DeleteBackupWithResponse(ctx, id)
	if res == nil {
		return errors.NewWrap(err, "error starting backup")
	} else if !checkSuccess(res.HTTPResponse, err) {
		return errors.NewWrap(err, "error starting backup: "+string(res.Body))
	}
	return nil
}

func (c *Client) ListBackups(ctx contexts.Context) ([]spec.Backup, *errors.Error) {
	res, err := c.client.ListBackupsWithResponse(ctx)
	if res == nil {
		return nil, errors.NewWrap(err, "error listing backups")
	} else if checkSuccess(res.HTTPResponse, err) {
		backups := res.JSON200
		return *backups, nil
	} else {
		return nil, errors.NewWrap(err, "error listing backups: "+string(res.Body))
	}
}

func (c *Client) StartRestore(ctx contexts.Context, backupId string, req spec.ContainersConfig) (*spec.Restore, *errors.Error) {
	res, err := c.client.StartRestoreWithResponse(ctx, backupId, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting restore")
	} else if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error starting restore: "+string(res.Body))
	}
}

func (c *Client) GetRestore(ctx contexts.Context, id string) (*spec.Restore, *errors.Error) {
	res, err := c.client.GetRestoreWithResponse(ctx, id)
	if res == nil {
		return nil, errors.NewWrap(err, "error finding restore")
	} else if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error finding restore: "+id+": "+string(res.Body))
	}
}

func (c *Client) RotateBackups(ctx contexts.Context, req spec.RotateInput) (*spec.Rotate, *errors.Error) {
	res, err := c.client.RotateBackupsWithResponse(ctx, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting rotate")
	} else if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error starting rotate: "+string(res.Body))
	}
}

func (c *Client) Version(ctx contexts.Context) (*spec.Version, *errors.Error) {
	res, err := c.client.GetVersionWithResponse(ctx)
	if res == nil {
		return nil, errors.NewWrap(err, "error getting version")
	} else if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error getting version:"+string(res.Body))
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
