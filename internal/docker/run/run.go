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
		return &MissingTemplate{Message: "No config for " + template.Alias}
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

func runContainer(ctx contexts.Context, client client.DockerClient, rendered types.ContainerTemplate) error {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", encodings.ToJsonIE(rendered))
		return nil
	}
	logging.Debug(ctx, "Executing:", encodings.ToJsonIE(rendered))

	id, err := client.RunContainer(rendered.Image, rendered.Cmd, rendered.Env, rendered.Volumes)
	if err == nil {
		err = waitAndLog(client, rendered.Alias, id)
	}

	removalErr := client.RemoveContainer(id)
	return errors.Join(err, removalErr)
}

func waitAndLog(client client.DockerClient, name string, id string) error {
	waitErr := client.WaitForContainer(name, id)
	logErr := client.ContainerLogs(id)
	return errors.Join(waitErr, logErr)
}
