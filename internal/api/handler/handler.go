package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/backup"
	"github.com/sbnarra/bckupr/internal/app/delete"
	"github.com/sbnarra/bckupr/internal/app/list"
	"github.com/sbnarra/bckupr/internal/app/restore"
	"github.com/sbnarra/bckupr/internal/app/rotate"
	"github.com/sbnarra/bckupr/internal/app/version"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/meta/reader"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/tracker"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Handler struct {
	Context              contexts.Context
	Templates            containers.Templates
	NotificationSettings *notifications.NotificationSettings
}

func (s Handler) newContext(c *gin.Context) contexts.Context {
	ctx := contexts.Copy(c, s.Context)
	ctx.Name = c.Request.URL.Path
	dryRunH := c.Request.Header.Get("Dry-Run")
	if dryRunH != "" {
		if dryRun, err := strconv.ParseBool(dryRunH); err == nil {
			ctx.DryRun = dryRun
		} else {
			ctx.DryRun = true
		}
	}
	debugH := c.Request.Header.Get("Debug")
	if debugH != "" {
		if debug, err := strconv.ParseBool(debugH); err == nil {
			ctx.Debug = debug
		} else {
			ctx.Debug = true
		}
	}
	return ctx
}

func (s Handler) StartBackup(c *gin.Context) {
	s.StartBackupWithId(c, "")
}

func (s Handler) StartBackupWithId(c *gin.Context, id string) {
	ctx := s.newContext(c)
	payload := spec.TaskInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if backup, runner, err := backup.Start(ctx, id, payload, s.Templates, s.NotificationSettings); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		go func() {
			if err := runner.Wait(); err != nil {
				logging.CheckError(ctx, err, "backup failure")
			}
		}()
		onSuccess(c, http.StatusOK, backup)
	}
}

func (s Handler) StartRestore(c *gin.Context, id string) {
	ctx := s.newContext(c)
	payload := spec.TaskInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(spec.RestoreStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if task, runner, err := restore.Start(ctx, id, payload, s.Templates, s.NotificationSettings); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		go func() {
			if err := runner.Wait(); err != nil {
				logging.CheckError(ctx, err, "restore failure")
			}
		}()
		onSuccess(c, http.StatusOK, task)
	}
}

func (s Handler) DeleteBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if err := delete.Delete(ctx, id); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		onSuccess(c, http.StatusOK, nil)
	}
}

func (s Handler) ListBackups(c *gin.Context) {
	ctx := s.newContext(c)
	if backups, err := list.ListBackups(ctx); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		onSuccess(c, http.StatusOK, backups)
	}
}

func (s Handler) GetBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if reader, err := reader.Load(ctx); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else if backup := reader.Get(id); backup == nil {
		onError(ctx, c, http.StatusNotFound, errors.New(id+" not found"))
	} else {
		onSuccess(c, http.StatusOK, backup)
	}
}

func (s Handler) GetRestore(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if restore, err := tracker.Get[spec.Restore]("restore", id); err != nil {
		onError(ctx, c, http.StatusNotFound, err)
	} else {
		onSuccess(c, http.StatusOK, restore)
	}
}

func (h Handler) StartRotate(c *gin.Context) {
	ctx := h.newContext(c)
	payload := spec.RotateInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if rotate, runner, err := rotate.Rotate(ctx, payload); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		go func() {
			if err := runner.Wait(); err != nil {
				logging.CheckError(ctx, err, "rotate failure")
			}
		}()
		onSuccess(c, http.StatusOK, rotate)
	}
}

func (h Handler) GetRotate(c *gin.Context) {
	ctx := h.newContext(c)
	if rotate, err := tracker.Get[spec.Rotate]("rotate", ""); err != nil {
		onError(ctx, c, http.StatusNotFound, err)
	} else {
		onSuccess(c, http.StatusOK, rotate)
	}
}

func (s Handler) GetVersion(c *gin.Context) {
	ctx := s.newContext(c)
	version := version.Version(ctx)
	onSuccess(c, http.StatusOK, version)
}

func onSuccess(c *gin.Context, status int, response any) {
	if response != nil {
		c.JSON(status, response)
	} else {
		c.Status(status)
	}
}

func onError(ctx contexts.Context, c *gin.Context, status int, err *errors.Error) {
	c.JSON(status, spec.Error{
		Error: err.Error(),
	})
	logging.CheckError(ctx, err, "Status", status)
}
