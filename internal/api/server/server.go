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
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type handler struct {
	ctx        contexts.Context
	containers containers.Templates
}

func (s handler) newContext(c *gin.Context) contexts.Context {
	ctx := contexts.Copy(c, s.ctx)
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

func (s handler) TriggerBackupWithId(c *gin.Context, id string) {
	ctx := s.newContext(c)

	payload := spec.BackupTrigger{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
		logging.CheckError(ctx, errors.Wrap(err, "failed to bind request"))
	} else if err := payload.WithDefaults(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		if backup, err := backup.CreateBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			logging.CheckError(ctx, err)
		} else {
			c.JSON(http.StatusOK, backup)
		}
	}
}

func (s handler) TriggerBackup(c *gin.Context) {
	s.TriggerBackupWithId(c, "")
}

func (s handler) TriggerRestore(c *gin.Context, id string) {
	ctx := s.newContext(c)
	payload := spec.RestoreTrigger{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
		logging.CheckError(ctx, errors.Wrap(err, "failed to bind request"))
	} else if err := payload.WithDefaults(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		if task, err := restore.RestoreBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, task)
		}
	}
}

func (s handler) DeleteBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if err := delete.DeleteBackup(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s handler) ListBackups(c *gin.Context) {
	ctx := s.newContext(c)
	if backups, err := list.ListBackups(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		c.JSON(http.StatusOK, backups)
	}
}

func (s handler) GetBackup(c *gin.Context, id string) {
	ctx := s.newContext(c)
	if reader, err := meta.NewReader(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		backup := reader.Get(id)
		if backup != nil {
			c.JSON(http.StatusOK, backup)
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}

func (s handler) RotateBackups(c *gin.Context) {
	ctx := s.newContext(c)
	payload := spec.RotateTrigger{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
		logging.CheckError(ctx, errors.Wrap(err, "failed to bind request"))
	} else if err := payload.WithDefaults(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		logging.CheckError(ctx, err)
	} else {
		if err := rotate.Rotate(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			logging.CheckError(ctx, err)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (s handler) GetVersion(c *gin.Context) {
	ctx := s.newContext(c)
	version := version.Version(ctx)
	c.JSON(http.StatusOK, version)
}
