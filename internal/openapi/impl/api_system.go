package impl

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/openapi/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type systemApi struct {
	ctx contexts.Context
}

func NewSystemAPI(ctx contexts.Context) spec.SystemAPI {
	return systemApi{
		ctx: ctx,
	}
}

func (api systemApi) GetVersion(c *gin.Context) {
	var created time.Time

	createdEnv := os.Getenv("CREATED")
	if createdEnv == "" {
		created = time.Now()
	} else {
		var err error
		created, err = time.Parse("", createdEnv)
		logging.CheckError(api.ctx, errors.Wrap(err, "failed to parse created date/time: "+createdEnv))
	}

	c.JSON(http.StatusOK, spec.Version{
		Version: os.Getenv("VERSION"),
		Created: created,
	})
}
