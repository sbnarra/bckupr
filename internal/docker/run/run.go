package run

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RunContainer(ctx contexts.Context, client client.DockerClient, meta CommonEnv, template types.ContainerTemplate, waitLogCleanup bool) (string, *errors.Error) {
	if len(template.Image) == 0 || len(template.Cmd) == 0 {
		return "", errors.Wrap(MisconfiguredTemplate, encodings.ToJsonIE(template))
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

func runContainer(ctx contexts.Context, client client.DockerClient, template types.ContainerTemplate, waitLogCleanup bool) (string, *errors.Error) {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", encodings.ToJsonIE(template))
		return "", nil
	}
	logging.Debug(ctx, "Executing:", encodings.ToJsonIE(template))

	id, err := client.RunContainer(ctx, template.Image, template.Cmd, template.Env, template.Volumes, template.Labels)

	// using new context as we always want to log/remove even if cancelled/interrupted
	ctx = contexts.NonCancallable(ctx)

	if !waitLogCleanup {
		return id, err
	} else if err == nil {
		err = WaitThenLog(ctx, client, id)
	}

	removalErr := client.RemoveContainer(ctx, id)
	return id, errors.Join(err, removalErr)
}

func WaitThenLog(ctx contexts.Context, client client.DockerClient, id string) *errors.Error {
	waitErr := client.WaitForContainer(ctx, id)
	logs, logErr := client.ContainerLogs(ctx, id)

	logCtx := ctx
	logCtx.Name = id
	logging.Debug(logCtx, logs)
	return errors.Join(waitErr, logErr)
}
