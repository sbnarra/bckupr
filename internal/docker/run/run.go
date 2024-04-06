package run

import (
	"errors"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RunContainer(ctx contexts.Context, client client.DockerClient, meta CommonEnv, template types.ContainerTemplate) error {
	if len(template.Image) == 0 || len(template.Cmd) == 0 {
		return &MisconfiguredTemplate{Message: "Misconfigured template: " + encodings.ToJsonIE(template)}
	}

	copy := template

	copy.Env = append(copy.Env,
		"VOLUME_NAME="+meta.VolumeName,
		"BACKUP_ID="+meta.BackupId,
		"FILE_EXT="+meta.FileExt,
		"BACKUP_DIR=/backup",
		"DATA_DIR=/data",
	)
	return runContainer(ctx, client, copy)
}

func runContainer(ctx contexts.Context, client client.DockerClient, template types.ContainerTemplate) error {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", encodings.ToJsonIE(template))
		return nil
	}
	logging.Debug(ctx, "Executing:", encodings.ToJsonIE(template))

	id, err := client.RunContainer(template.Image, template.Cmd, template.Env, template.Volumes)
	if err == nil {
		err = waitAndLog(ctx, client, id)
	}

	removalErr := client.RemoveContainer(id)
	return errors.Join(err, removalErr)
}

func waitAndLog(ctx contexts.Context, client client.DockerClient, id string) error {
	waitErr := client.WaitForContainer(ctx, id)
	logErr := client.ContainerLogs(id)
	return errors.Join(waitErr, logErr)
}
