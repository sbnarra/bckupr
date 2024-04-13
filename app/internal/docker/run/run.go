package run

import (
	"errors"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/encodings"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
)

func RunContainer(ctx contexts.Context, client client.DockerClient, meta CommonEnv, template types.ContainerTemplate, waitLogCleanup bool) (string, error) {
	if len(template.Image) == 0 || len(template.Cmd) == 0 {
		return "", &MisconfiguredTemplate{Message: "Misconfigured template: " + encodings.ToJsonIE(template)}
	}

	copy := template

	copy.Env = append(copy.Env,
		"VOLUME_NAME="+meta.VolumeName,
		"BACKUP_ID="+meta.BackupId,
		"FILE_EXT="+meta.FileExt,
		"BACKUP_DIR=/backup",
		"DATA_DIR=/data",
	)
	return runContainer(ctx, client, copy, waitLogCleanup)
}

func runContainer(ctx contexts.Context, client client.DockerClient, template types.ContainerTemplate, waitLogCleanup bool) (string, error) {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", encodings.ToJsonIE(template))
		return "", nil
	}
	logging.Debug(ctx, "Executing:", encodings.ToJsonIE(template))

	id, err := client.RunContainer(template.Image, template.Cmd, template.Env, template.Volumes, template.Labels)

	if !waitLogCleanup {
		return id, err
	} else if err == nil {
		err = WaitThenLog(ctx, client, id)
	}

	removalErr := client.RemoveContainer(id)
	return id, errors.Join(err, removalErr)
}

func WaitThenLog(ctx contexts.Context, client client.DockerClient, id string) error {
	waitErr := client.WaitForContainer(ctx, id)
	logs, logErr := client.ContainerLogs(id)

	logCtx := ctx
	logCtx.Name = id
	logging.Debug(logCtx, logs)
	return errors.Join(waitErr, logErr)
}
