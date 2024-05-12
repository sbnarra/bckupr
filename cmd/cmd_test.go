package cmd

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func TestCmdE2E(t *testing.T) {
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})

	e2e.PrepareIntegrationTest(t)

	go func() {
		os.Args = []string{"", "daemon",
			"--" + keys.HostBackupDir.CliId, e2e.BackupDir,
			"--" + keys.ContainerBackupDir.CliId, e2e.BackupDir}
		Cmd.Execute()
	}()
	time.Sleep(time.Second * 3)
	defer func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	id := time.Now().Format("200601021504") + "-cmd"
	e2e.RunE2E(t,
		func() *errors.E {

			os.Args = []string{"", "backup", "--" + keys.NoDryRun.CliId, "--" + keys.BackupId.CliId, id}
			return errors.Wrap(Cmd.Execute(), "backup failed")
		},
		func() *errors.E {
			os.Args = []string{"", "restore", "--" + keys.NoDryRun.CliId, "--" + keys.BackupId.CliId, id}
			return errors.Wrap(Cmd.Execute(), "backup failed")
		},
		func() *errors.E {
			os.Args = []string{"", "delete", "--" + keys.BackupId.CliId, id}
			return errors.Wrap(Cmd.Execute(), "backup failed")
		})
}
