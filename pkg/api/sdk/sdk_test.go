package sdk_test

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/internal/web/server"
	"github.com/sbnarra/bckupr/pkg/api/sdk"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func TestSdkE2E(t *testing.T) {
	ctx := e2e.PrepareIntegrationTest(t)

	config := e2e.NewServerConfig()
	if containers, err := containers.LoadTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig); err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	} else {
		s := server.New(ctx, config, containers, &notifications.NotificationSettings{})
		go func() {
			if err := s.Listen(ctx); err != nil {
				logging.CheckWarn(ctx, err)
			}
		}()
		defer s.Server.Shutdown(ctx)
	}

	time.Sleep(2 * time.Second)

	sdk, err := sdk.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string))
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	id := time.Now().Format("20060102_1504") + "-sdk"
	taskInput := spec.TaskInput{}

	e2e.RunE2E(t,
		func() *errors.Error {
			if _, err := sdk.StartBackupWithId(ctx, id, taskInput); err != nil {
				return err
			} else {
				for backup, err := sdk.GetBackup(ctx, id); err == nil; {
					logging.Info(ctx, "backup:", backup.Status)
					if backup.Status == spec.StatusCompleted {
						return nil
					} else if backup.Status == spec.StatusError {
						return errors.New(*backup.Error)
					}
					time.Sleep(time.Second * 2)
					backup, err = sdk.GetBackup(ctx, id)
				}
				return err
			}
		},
		func() *errors.Error {
			if _, err := sdk.StartRestore(ctx, id, taskInput); err != nil {
				return err
			} else {
				for restore, err := sdk.GetRestore(ctx, id); err == nil; {
					logging.Info(ctx, "restore:", restore.Status)
					if restore.Status == spec.StatusCompleted {
						return nil
					} else if restore.Status == spec.StatusError {
						return errors.New(*restore.Error)
					}

					time.Sleep(time.Second * 2)
					restore, err = sdk.GetRestore(ctx, id)
				}
				return err
			}
		},
		func() *errors.Error {
			return sdk.DeleteBackup(ctx, id)
		})
}
