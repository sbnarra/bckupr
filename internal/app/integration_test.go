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
	"github.com/sbnarra/bckupr/internal/tasks/async"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func TestAppE2E(t *testing.T) {

	ctx := e2e.PrepareIntegrationTest(t)
	id := time.Now().Format("20060102_1504") + "-internal"

	config := e2e.NewServerConfig()
	containers, err := containers.ContainerTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig)
	if err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	}

	notificationSettings := &notifications.NotificationSettings{}
	e2e.RunE2E(t,
		func() *errors.Error {
			payload := spec.ContainersConfig{}
			if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
				return err
			} else {
				if backup, err := backup.Start(ctx, "", payload, containers, notificationSettings); err != nil {
					return err
				} else {
					id = backup.Id
					if async, err := async.Current("backup", id); err != nil {
						return err
					} else {
						return async.Runner.Wait()
					}
				}
			}
		},
		func() *errors.Error {
			payload := spec.ContainersConfig{}
			if err := payload.WithDefaults(spec.BackupStopModes); err != nil {
				return err
			} else {
				if _, err := restore.Start(ctx, id, payload, containers, notificationSettings); err != nil {
					return err
				} else {
					if async, err := async.Current("restore", id); err != nil {
						return err
					} else {
						return async.Runner.Wait()
					}
				}
				return err
			}
		},
		func() *errors.Error {
			return delete.Delete(ctx, id)
		})
}
