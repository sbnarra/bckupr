package impl

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/openapi/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (api backupAPI) createContext(c *gin.Context) contexts.Context {
	ctx := contexts.Copy(c, api.ctx)

	debug := c.Request.Header.Get("debug")
	ctx.Debug, _ = strconv.ParseBool(debug)

	dryRun := c.Request.Header.Get("dryRun")
	ctx.DryRun, _ = strconv.ParseBool(dryRun)

	return ctx
}

func (api backupAPI) handleWithPayload(c *gin.Context, input any, exec func(contexts.Context) *errors.Error) {
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, spec.Error{
			Error: err.Error(),
		})
		return
	}
	api.handle(c, exec)
}

func (api backupAPI) handle(c *gin.Context, exec func(contexts.Context) *errors.Error) {
	ctx := api.createContext(c)
	if err := exec(ctx); err == nil {
		c.Status(http.StatusOK)
	} else {
		internalError(c, ctx, err)
	}
}

func internalError(c *gin.Context, ctx contexts.Context, err *errors.Error) {
	logging.CheckError(ctx, err, "error processing request")
	c.JSON(http.StatusInternalServerError, spec.Error{
		Error: err.Error(),
	})
}

func notFound(c *gin.Context, ctx contexts.Context, msg string) {
	c.JSON(http.StatusNotFound, spec.Error{
		Error: msg,
	})
}
