package app_test

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/backup"
	"github.com/sbnarra/bckupr/internal/app/delete"
	"github.com/sbnarra/bckupr/internal/app/restore"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func TestAppE2E(t *testing.T) {

	ctx := e2e.PrepareIntegrationTest(t)
	id := time.Now().Format("20060102_150400") + "-app"

	config := e2e.NewServerConfig()
	containers, err := containers.LoadTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig)
	if err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	}

	tr := true
	notificationSettings := &notifications.NotificationSettings{}
	e2e.RunE2E(t,
		func() *errors.E {
			payload := spec.TaskInput{
				NoDryRun: &tr,
			}
			if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
				return err
			} else {
				if _, runner, err := backup.Start(ctx, id, e2e.DockerHosts, e2e.BackupDir, e2e.BackupDir, payload, containers, notificationSettings); err != nil {
					return err
				} else {
					return runner.Wait()
				}
			}
		},
		func() *errors.E {
			payload := spec.TaskInput{
				NoDryRun: &tr,
			}
			if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
				return err
			} else {
				if _, runner, err := restore.Start(ctx, id, e2e.DockerHosts, e2e.BackupDir, e2e.BackupDir, payload, containers, notificationSettings); err != nil {
					return err
				} else {
					return runner.Wait()
				}
			}
		},
		func() *errors.E {
			return delete.Delete(ctx, id, e2e.BackupDir)
		})
}
