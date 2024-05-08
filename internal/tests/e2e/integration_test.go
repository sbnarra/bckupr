package e2e

import (
	"fmt"
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
			createBackup := serverSpec.NewTriggerBackup()
			task, backup, err := backup.CreateBackup(ctx, "", createBackup, containers)
			fmt.Println(task)
			id = backup.Id
			return err
		},
		func() *errors.Error {
			restoreBackup := serverSpec.NewTriggerRestore()
			task, err := restore.RestoreBackup(ctx, id, restoreBackup, containers)
			logging.Info(ctx, task)
			return err
		},
		func() *errors.Error {
			return delete.DeleteBackup(ctx, id)
		})
}

func TestE2EExternal(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	config := config.New()
	if containers, err := containers.ContainerTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig); err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	} else {
		if err := server.Start(ctx, config, nil, containers); err != nil {
			t.Fatalf("failed to start", err)
		}
	}

	time.Sleep(2 * time.Second)

	client, err := client.New(keys.DaemonAddr.Default.(string))
	if err != nil {
		t.Fatalf("error creating client: %w", err)
	}
	id := time.Now().Format("20060102_1504") + "-client"

	e2e(t,
		func() *errors.Error {
			req := clientSpec.BackupTrigger{}
			task, backup, err := client.TriggerBackupUsingId(ctx, id, req)
			logging.Info(ctx, task, backup, err)
			return err
		},
		func() *errors.Error {
			req := clientSpec.RestoreTrigger{}
			task, err := client.TriggerRestore(ctx, id, req)
			logging.Info(ctx, task, err)
			return err
		},
		func() *errors.Error {
			return client.DeleteBackup(ctx, id)
		})
}
