package e2e

import (
	"context"
	"testing"

	ctx "github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/web/server"
)

var DockerHosts = []string{"unix:///var/run/docker.sock"}
var BackupDir = "/tmp/backups"

func createContext(t *testing.T) context.Context {
	debug := true
	return ctx.Using(context.Background(), t.Name(), debug, 1)
}

func NewServerConfig() server.Config {
	return server.Config{
		DockerHosts:        DockerHosts,
		HostBackupDir:      BackupDir,
		ContainerBackupDir: BackupDir,

		LocalContainersConfig:   keys.LocalContainersConfig.EnvString(),
		OffsiteContainersConfig: keys.OffsiteContainersConfig.EnvString(),
		TcpAddr:                 keys.TcpAddr.EnvString(),

		NotificationSettings: &notifications.NotificationSettings{},
	}
}
