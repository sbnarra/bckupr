package client_test

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/api/server"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/async"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/client"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func TestClientE2E(t *testing.T) {
	ctx := e2e.PrepareIntegrationTest(t)

	config := e2e.NewServerConfig()
	if containers, err := containers.ContainerTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig); err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	} else {
		s := server.New(ctx, config, containers, &notifications.NotificationSettings{})
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
	containersConfig := spec.ContainersConfig{}

	e2e.RunE2E(t,
		func() *errors.Error {
			if backup, err := client.StartBackupWithId(ctx, id, containersConfig); err != nil {
				return err
			} else {
				id = backup.Id
				if async, err := async.Current("backup", id); err != nil {
					return err
				} else {
					return async.Runner.Wait()
				}
			}
		},
		func() *errors.Error {
			if _, err := client.StartRestore(ctx, id, containersConfig); err != nil {
				return err
			} else {
				if async, err := async.Current("restore", id); err != nil {
					return err
				} else {
					return async.Runner.Wait()
				}
			}
		},
		func() *errors.Error {
			return client.DeleteBackup(ctx, id)
		})
}
