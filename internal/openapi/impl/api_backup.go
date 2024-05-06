package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/openapi/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

type backupAPI struct {
	ctx        contexts.Context
	containers types.ContainerTemplates
}

func NewBackupAPI(ctx contexts.Context, containers types.ContainerTemplates) spec.BackupAPI {
	return backupAPI{
		ctx:        ctx,
		containers: containers,
	}
}

func (api backupAPI) CreateBackup(c *gin.Context) {
	input := types.DefaultCreateBackupRequest()
	api.handleWithPayload(c, input, func(ctx contexts.Context) *errors.Error {
		_, err := app.CreateBackup(ctx, input, api.containers)
		return err
	})
}

func (api backupAPI) GetBackup(c *gin.Context) {
	id := c.Param("id")
	ctx := api.createContext(c)
	if mw, err := meta.NewReader(ctx); err != nil {
		internalError(c, ctx, err)
	} else if backup := mw.Get(id); backup != nil {
		c.JSON(http.StatusOK, backup)
	} else {
		notFound(c, ctx, "backup not found: "+id)
	}
}

func (api backupAPI) DeleteBackup(c *gin.Context) {
	api.handle(c, func(ctx contexts.Context) *errors.Error {
		return app.DeleteBackup(ctx, c.Param("id"))
	})
}

func (api backupAPI) RestoreBackup(c *gin.Context) {
	input := types.DefaultRestoreBackupRequest()
	api.handleWithPayload(c, input, func(ctx contexts.Context) *errors.Error {
		return app.RestoreBackup(ctx, input, api.containers)
	})
}
