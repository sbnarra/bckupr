package sdk

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

type Sdk struct {
	client *spec.ClientWithResponses
}

func New(ctx context.Context, protocol string, server string) (*Sdk, *errors.E) {
	url := protocol + "://" + server
	if http, err := spec.NewClientWithResponses(url, func(c *spec.Client) error {
		c.RequestEditors = []spec.RequestEditorFn{
			func(ct context.Context, req *http.Request) error {
				req.Header.Add("User-Agent", "bckupr-sdk/"+os.Getenv("VERSION"))
				req.Header.Add("Debug", strconv.FormatBool(contexts.Debug(ctx)))
				return nil
			},
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create api client")
	} else {
		return &Sdk{http}, nil
	}
}

func (s *Sdk) StartBackup(ctx context.Context, req spec.TaskInput) (*spec.Backup, *errors.E) {
	res, err := s.client.StartBackupWithResponse(ctx, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error starting backup: "+string(res.Body))
	}
}

func (s *Sdk) StartBackupWithId(ctx context.Context, id string, req spec.TaskInput) (*spec.Backup, *errors.E) {
	res, err := s.client.StartBackupWithIdWithResponse(ctx, id, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error starting backup:"+id+": "+string(res.Body))
	}
}

func (s *Sdk) GetBackup(ctx context.Context, id string) (*spec.Backup, *errors.E) {
	res, err := s.client.GetBackupWithResponse(ctx, id)
	if res == nil {
		return nil, errors.NewWrap(err, "error finding backup")
	} else if checkSuccess(res.HTTPResponse, err) {
		backup := res.JSON200
		return backup, nil
	} else {
		return nil, errors.NewWrap(err, "error finding backup: "+id+": "+string(res.Body))
	}
}

func (s *Sdk) DeleteBackup(ctx context.Context, id string) *errors.E {
	res, err := s.client.DeleteBackupWithResponse(ctx, id)
	if res == nil {
		return errors.NewWrap(err, "error starting backup")
	} else if !checkSuccess(res.HTTPResponse, err) {
		return errors.NewWrap(err, "error starting backup: "+string(res.Body))
	}
	return nil
}

func (s *Sdk) ListBackups(ctx context.Context) ([]spec.Backup, *errors.E) {
	res, err := s.client.ListBackupsWithResponse(ctx)
	if res == nil {
		return nil, errors.NewWrap(err, "error listing backups")
	} else if checkSuccess(res.HTTPResponse, err) {
		backups := res.JSON200
		return *backups, nil
	} else {
		return nil, errors.NewWrap(err, "error listing backups: "+string(res.Body))
	}
}

func (s *Sdk) StartRestore(ctx context.Context, backupId string, req spec.TaskInput) (*spec.Restore, *errors.E) {
	res, err := s.client.StartRestoreWithResponse(ctx, backupId, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting restore")
	} else if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error starting restore: "+string(res.Body))
	}
}

func (s *Sdk) GetRestore(ctx context.Context, id string) (*spec.Restore, *errors.E) {
	res, err := s.client.GetRestoreWithResponse(ctx, id)
	if res == nil {
		return nil, errors.NewWrap(err, "error finding restore")
	} else if checkSuccess(res.HTTPResponse, err) {
		restore := res.JSON200
		return restore, nil
	} else {
		return nil, errors.NewWrap(err, "error finding restore: "+id+": "+string(res.Body))
	}
}

func (s *Sdk) StartRotate(ctx context.Context, req spec.RotateInput) (*spec.Rotate, *errors.E) {
	res, err := s.client.StartRotateWithResponse(ctx, req)
	if res == nil {
		return nil, errors.NewWrap(err, "error starting rotate")
	} else if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error starting rotate: "+string(res.Body))
	}
}

func (s *Sdk) GetRotate(ctx context.Context) (*spec.Rotate, *errors.E) {
	res, err := s.client.GetRotateWithResponse(ctx)
	if res == nil {
		return nil, errors.NewWrap(err, "error finding rotate")
	} else if checkSuccess(res.HTTPResponse, err) {
		task := res.JSON200
		return task, nil
	} else {
		return nil, errors.NewWrap(err, "error finding rotate: "+string(res.Body))
	}
}

func (s *Sdk) Version(ctx context.Context) (*spec.Version, *errors.E) {
	res, err := s.client.GetVersionWithResponse(ctx)
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
