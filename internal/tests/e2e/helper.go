package e2e

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/web/server"
)

func createContext(t *testing.T) contexts.Context {
	debug := true
	dryRun := false
	backupDir := "/tmp/backups"
	dockerHosts := []string{"unix:///var/run/docker.sock"}

	os.Setenv(keys.Debug.EnvId(), strconv.FormatBool(debug))
	os.Setenv(keys.DryRun.EnvId(), strconv.FormatBool(dryRun))
	os.Setenv(keys.HostBackupDir.EnvId(), backupDir)

	return contexts.Create(context.Background(), t.Name(), 1, backupDir, backupDir, dockerHosts, contexts.Debug(debug), contexts.DryRun(dryRun))
}

func NewServerConfig() server.Config {
	return server.Config{
		DockerHosts:             keys.DockerHosts.EnvStringSlice(),
		HostBackupDir:           keys.HostBackupDir.EnvString(),
		LocalContainersConfig:   keys.LocalContainersConfig.EnvString(),
		OffsiteContainersConfig: keys.OffsiteContainersConfig.EnvString(),
		TcpAddr:                 keys.TcpAddr.EnvString(),
	}
}
