package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/api/config"
	"github.com/sbnarra/bckupr/internal/api/server"
	serverSpec "github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/backup"
	"github.com/sbnarra/bckupr/internal/app/delete"
	"github.com/sbnarra/bckupr/internal/app/restore"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/client"
	clientSpec "github.com/sbnarra/bckupr/pkg/api/spec"
)

// move into internal/app package
func TestE2EInternal(t *testing.T) {

	ctx := prepareIntegrationTest(t)
	id := time.Now().Format("20060102_1504") + "-internal"

	config := config.New()
	containers, err := containers.ContainerTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig)
	if err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	}

	e2e(t,
		func() *errors.Error {
			payload := serverSpec.ContainersConfig{}
			if err := payload.WithDefaults(serverSpec.BackupStopModes); err != nil {
				return err
			} else {
				backup, err := backup.CreateBackup(ctx, "", payload, containers)
				id = backup.Id
				return err
			}
		},
		func() *errors.Error {
			payload := serverSpec.ContainersConfig{}
			if err := payload.WithDefaults(serverSpec.BackupStopModes); err != nil {
				return err
			} else {
				task, err := restore.RestoreBackup(ctx, id, payload, containers)
				logging.Info(ctx, task)
				return err
			}
		},
		func() *errors.Error {
			return delete.DeleteBackup(ctx, id)
		})
}

// move into pkg/client package
func TestE2EExternal(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	config := config.New()
	if containers, err := containers.ContainerTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig); err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	} else {
		s := server.New(ctx, config, containers)
		go func() {
			if err := s.Listen(ctx); err != nil {
				logging.CheckWarn(ctx, err)
			}
		}()
		defer s.Server.Close()
	}

	time.Sleep(2 * time.Second)

	client, err := client.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string))
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	id := time.Now().Format("20060102_1504") + "-client"

	e2e(t,
		func() *errors.Error {
			req := clientSpec.ContainersConfig{}
			backup, err := client.TriggerBackupUsingId(ctx, id, req)
			logging.Info(ctx, backup, err)
			return err
		},
		func() *errors.Error {
			req := clientSpec.ContainersConfig{}
			task, err := client.TriggerRestore(ctx, id, req)
			logging.Info(ctx, task, err)
			return err
		},
		func() *errors.Error {
			return client.DeleteBackup(ctx, id)
		})
}
