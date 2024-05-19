package run

import (
	"context"

	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func RunContainer(ctx context.Context, client client.DockerClient, meta CommonEnv, template containers.Template, waitLogCleanup bool) (string, *errors.E) {
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

	if copy.Labels == nil {
		copy.Labels = map[string]string{}
	}
	copy.Labels["managedby"] = "bckupr"
	copy.Labels["bckupr.backupid"] = meta.BackupId
	copy.Labels["bckupr.volume"] = meta.VolumeName

	return runContainer(ctx, client, copy, waitLogCleanup)
}

func runContainer(ctx context.Context, client client.DockerClient, template containers.Template, waitLogCleanup bool) (string, *errors.E) {
	id, runErr := client.RunContainer(ctx, template.Image, template.Cmd, template.Env, template.Volumes, template.Labels)
	if runErr != nil || !waitLogCleanup {
		return id, runErr
	} else if id == "dry_run" {
		return id, nil
	}

	ctx = context.WithoutCancel(ctx)
	waitLogErr := WaitThenLog(ctx, client, id)
	removalErr := client.RemoveContainer(ctx, id)
	return id, errors.Join(waitLogErr, removalErr)
}

func WaitThenLog(ctx context.Context, client client.DockerClient, id string) *errors.E {
	waitErr := client.WaitForContainer(ctx, id)
	logs, logErr := client.ContainerLogs(ctx, id)

	name := contexts.Name(ctx)
	logCtx := contexts.WithName(ctx, name+"/"+id[:7]+":OUT")
	logging.Debug(logCtx, logs.Out)
	logCtx = contexts.WithName(ctx, name+"/"+id[:7]+":ERR")
	logging.Debug(logCtx, logs.Err)

	if waitErr != nil {
		return errors.Wrap(waitErr, logs.Err)
	}
	return logErr
}
