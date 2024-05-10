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
	"github.com/sbnarra/bckupr/internal/tasks/async"
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
				if backup, err := backup.Start(ctx, "", payload, containers); err != nil {
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
			payload := serverSpec.ContainersConfig{}
			if err := payload.WithDefaults(serverSpec.BackupStopModes); err != nil {
				return err
			} else {
				if _, err := restore.Start(ctx, id, payload, containers); err != nil {
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
			if backup, err := client.StartBackupWithId(ctx, id, req); err != nil {
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
			req := clientSpec.ContainersConfig{}
			if _, err := client.StartRestore(ctx, id, req); err != nil {
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
