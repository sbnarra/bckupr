package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/backup"
	"github.com/sbnarra/bckupr/internal/app/delete"
	"github.com/sbnarra/bckupr/internal/app/list"
	"github.com/sbnarra/bckupr/internal/app/restore"
	"github.com/sbnarra/bckupr/internal/app/rotate"
	"github.com/sbnarra/bckupr/internal/app/version"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

type handler struct {
	ctx        contexts.Context
	containers containers.Templates
}

func (s handler) TriggerBackupWithId(c *gin.Context, id string) {
	payload := spec.NewTriggerBackup()
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		ctx := contexts.Copy(c, s.ctx)
		if task, backup, err := backup.CreateBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			fmt.Println("TODO: return type should include task and backup...or task should embed backup", task)
			c.JSON(http.StatusOK, backup)
		}
	}
}

func (s handler) TriggerBackup(c *gin.Context) {
	s.TriggerBackupWithId(c, "")
}

func (s handler) TriggerRestore(c *gin.Context, id string) {
	payload := spec.NewTriggerRestore()
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		ctx := contexts.Copy(c, s.ctx)
		if task, err := restore.RestoreBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, task)
		}
	}
}

func (s handler) DeleteBackup(c *gin.Context, id string) {
	ctx := contexts.Copy(c, s.ctx)
	if err := delete.DeleteBackup(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s handler) ListBackups(c *gin.Context) {
	ctx := contexts.Copy(c, s.ctx)
	if err := list.ListBackups(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s handler) GetBackup(c *gin.Context, id string) {}

func (s handler) RotateBackups(c *gin.Context) {
	payload := spec.NewRotateTrigger()
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		ctx := contexts.Copy(c, s.ctx)
		if err := rotate.Rotate(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (s handler) GetVersion(c *gin.Context) {
	ctx := contexts.Copy(c, s.ctx)
	version := version.Version(ctx)
	c.JSON(http.StatusOK, version)
}
