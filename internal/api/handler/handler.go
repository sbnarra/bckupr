package handler

import (
	"context"
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
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/meta/reader"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/tracker"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type handler struct {
	Context context.Context
	// move into internal/tasks/config
	DockerHosts          []string
	containerBackupDir   string
	hostBackupDir        string
	Templates            containers.Templates
	NotificationSettings *notifications.NotificationSettings
}

func New(
	ctx context.Context,
	dockerHosts []string,
	containerBackupDir string,
	hostBackupDir string,
	containers containers.Templates,
	notificationSettings *notifications.NotificationSettings,
) handler {
	return handler{
		Context:              ctx,
		DockerHosts:          dockerHosts,
		containerBackupDir:   containerBackupDir,
		hostBackupDir:        hostBackupDir,
		Templates:            containers,
		NotificationSettings: notificationSettings,
	}
}

func (h handler) newContext(c *gin.Context) context.Context {
	debug := false
	debugH := c.Request.Header.Get("Debug")
	if debugH != "" {
		if debugHB, err := strconv.ParseBool(debugH); err == nil {
			debug = debugHB
		}
	}
	threadLimit := contexts.ThreadLimit(c)
	return contexts.Using(c, c.Request.URL.Path, debug, threadLimit)
}

func (h handler) StartBackup(c *gin.Context) {
	h.StartBackupWithId(c, "")
}

func (h handler) StartBackupWithId(c *gin.Context, id string) {
	ctx := h.newContext(c)
	payload := spec.TaskInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if backup, runner, err := backup.Start(ctx, id, h.DockerHosts, h.hostBackupDir, h.containerBackupDir, payload, h.Templates, h.NotificationSettings); err != nil {
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

func (h handler) StartRestore(c *gin.Context, id string) {
	ctx := h.newContext(c)
	payload := spec.TaskInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(spec.RestoreStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if task, runner, err := restore.Start(ctx, id, h.DockerHosts, h.hostBackupDir, h.containerBackupDir, payload, h.Templates, h.NotificationSettings); err != nil {
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

func (h handler) DeleteBackup(c *gin.Context, id string) {
	ctx := h.newContext(c)
	if err := delete.Delete(ctx, id, h.containerBackupDir); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		onSuccess(c, http.StatusOK, nil)
	}
}

func (h handler) ListBackups(c *gin.Context) {
	ctx := h.newContext(c)
	if backups, err := list.ListBackups(ctx, h.containerBackupDir); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err)
	} else {
		onSuccess(c, http.StatusOK, backups)
	}
}

func (h handler) GetBackup(c *gin.Context, id string) {
	ctx := h.newContext(c)
	if backup, err := tracker.Get[spec.Backup]("backup", id); err != nil {
		if backup, err := h.GetBackupFromDisk(ctx, id); err != nil {
			onError(ctx, c, http.StatusNotFound, err)
		} else {
			onSuccess(c, http.StatusOK, backup)
		}
	} else {
		onSuccess(c, http.StatusOK, backup)
	}
}

func (h handler) GetRestore(c *gin.Context, id string) {
	ctx := h.newContext(c)
	if restore, err := tracker.Get[spec.Restore]("restore", id); err != nil {
		onError(ctx, c, http.StatusNotFound, err)
	} else {
		onSuccess(c, http.StatusOK, restore)
	}
}

func (h handler) GetBackupFromDisk(ctx context.Context, id string) (*spec.Backup, *errors.E) {
	if reader, err := reader.Load(ctx, h.containerBackupDir); err != nil {
		return nil, err
	} else if backup := reader.Get(id); backup == nil {
		return nil, errors.New(id + " not found")
	} else {
		return backup, nil
	}
}

func (h handler) StartRotate(c *gin.Context) {
	ctx := h.newContext(c)
	payload := spec.RotateInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, errors.Wrap(err, "error parsing request:"))
	} else if err := payload.WithDefaults(); err != nil {
		onError(ctx, c, http.StatusInternalServerError, errors.Wrap(err, "failed to load defaults:"))
	} else if rotate, runner, err := rotate.Rotate(ctx, payload, h.containerBackupDir); err != nil {
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

func (h handler) GetRotate(c *gin.Context) {
	ctx := h.newContext(c)
	if rotate, err := tracker.Get[spec.Rotate]("rotate", ""); err != nil {
		onError(ctx, c, http.StatusNotFound, err)
	} else {
		onSuccess(c, http.StatusOK, rotate)
	}
}

func (h handler) GetVersion(c *gin.Context) {
	ctx := h.newContext(c)
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

func onError(ctx context.Context, c *gin.Context, status int, err *errors.E) {
	c.JSON(status, spec.Error{
		Error: err.Error(),
	})
	logging.CheckError(ctx, err, "Status", status)
}
