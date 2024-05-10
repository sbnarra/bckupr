package server

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
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type handler struct {
	ctx        contexts.Context
	containers containers.Templates
}

func (s handler) newContext(c *gin.Context) contexts.Context {
	ctx := contexts.Copy(c, s.ctx)
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

func (s handler) StartBackup(c *gin.Context) {
	s.StartBackupWithId(c, "")
}

func (s handler) StartBackupWithId(c *gin.Context, id string) {
	ctx := s.newContext(c)
	payload := spec.ContainersConfig{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, "error parsing request:"+err.Error())
	} else if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, "failed to load defaults:"+err.Error())
	} else if backup, err := backup.Start(ctx, id, payload, s.containers); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else {
		onSuccess(c, http.StatusOK, backup)
	}
}

func (s handler) StartRestore(c *gin.Context, id string) {
	ctx := s.newContext(c)
	payload := spec.ContainersConfig{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, "error parsing request:"+err.Error())
	} else if err := payload.WithDefaults(spec.RestoreStopModes); err != nil {
		onError(ctx, c, http.StatusInternalServerError, "failed to load defaults:"+err.Error())
	} else if task, err := restore.Start(ctx, id, payload, s.containers); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else {
		onSuccess(c, http.StatusOK, task)
	}
}

func (s handler) DeleteBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if err := delete.Delete(ctx, id); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else {
		onSuccess(c, http.StatusOK, nil)
	}
}

func (s handler) ListBackups(c *gin.Context) {
	ctx := s.newContext(c)
	if backups, err := list.ListBackups(ctx); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else {
		onSuccess(c, http.StatusOK, backups)
	}
}

func (s handler) GetBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if reader, err := reader.Load(ctx); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else if backup := reader.Get(id); backup == nil {
		onError(ctx, c, http.StatusNotFound, id+" not found")
	} else {
		onSuccess(c, http.StatusOK, backup)
	}
}

func (s handler) GetRestore(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if latest := restore.Latest(); latest == nil {
		onError(ctx, c, http.StatusNotFound, id+" not found")
	} else {
		onSuccess(c, http.StatusOK, latest)
	}
}

func (s handler) RotateBackups(c *gin.Context) {
	ctx := s.newContext(c)
	payload := spec.RotateInput{}
	if err := c.BindJSON(&payload); err != nil {
		onError(ctx, c, http.StatusBadRequest, "error parsing request:"+err.Error())
	} else if err := payload.WithDefaults(); err != nil {
		onError(ctx, c, http.StatusInternalServerError, "failed to load defaults:"+err.Error())
	} else if err := rotate.Rotate(ctx, payload); err != nil {
		onError(ctx, c, http.StatusInternalServerError, err.Error())
	} else {
		onSuccess(c, http.StatusOK, nil)
	}
}

func (s handler) GetVersion(c *gin.Context) {
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

func onError(ctx contexts.Context, c *gin.Context, status int, err string) {
	c.JSON(status, spec.Error{
		Error: err,
	})
	logging.Error(ctx, err, status)
}
