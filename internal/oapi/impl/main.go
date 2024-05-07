package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/oapi/server"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

type Server struct {
	ctx        contexts.Context
	containers types.ContainerTemplates
}

func newServer(ctx contexts.Context, containers types.ContainerTemplates) *http.Server {
	router := gin.Default()
	server.RegisterHandlers(router, Server{ctx, containers})
	return &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
	}
	// log.Fatal(s.ListenAndServe())
}

func (s Server) TriggerBackupWithId(c *gin.Context, id string) {
	ctx := contexts.Copy(c, s.ctx)
	payload := server.NewTriggerBackup()
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		if backup, err := app.CreateBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, backup)
		}
	}
}

func (s Server) TriggerBackup(c *gin.Context) {
	s.TriggerBackupWithId(c, "")
}

func (s Server) TriggerRestore(c *gin.Context, id string) {
	ctx := contexts.Copy(c, s.ctx)
	payload := server.NewTriggerRestore()
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		if err := app.RestoreBackup(ctx, id, payload, s.containers); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (s Server) DeleteBackup(c *gin.Context, id string) {
	ctx := contexts.Copy(c, s.ctx)
	if err := app.DeleteBackup(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.Status(http.StatusOK)
	}
}
func (s Server) ListBackups(c *gin.Context) {}

func (s Server) GetBackup(c *gin.Context, id string) {}

func (s Server) GetVersion(c *gin.Context) {}
