package run

import (
	"errors"
	"strings"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RunContainer(ctx contexts.Context, client client.DockerClient, meta RunMeta, template types.ContainerTemplate) error {
	logging.Debug(ctx, "Template:", templateString(template))
	if rendered, err := renderTemplate(template, meta); err != nil {
		return err
	} else {
		return runContainer(ctx, client, rendered)
	}
}

func runContainer(ctx contexts.Context, client client.DockerClient, rendered types.ContainerTemplate) error {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", templateString(rendered))
		return nil
	}
	logging.Info(ctx, "Executing:", templateString(rendered))

	id, err := client.RunContainer(rendered.Image, rendered.Cmd, rendered.Env, rendered.Volumes)
	if err == nil {
		err = waitAndLog(ctx, client, rendered.Alias, id)
	}

	removalErr := client.RemoveContainer(id)
	return errors.Join(err, removalErr)
}

func renderTemplate(template types.ContainerTemplate, meta RunMeta) (types.ContainerTemplate, error) {
	if len(template.Image) == 0 || len(template.Cmd) == 0 {
		return types.ContainerTemplate{}, &MissingTemplate{Message: "No config for " + template.Alias}
	}

	formattedCmd := make([]string, 0)
	for _, part := range template.Cmd {
		part = strings.ReplaceAll(part, "{backup_id}", meta.BackupId)
		part = strings.ReplaceAll(part, "{name}", meta.VolumeName)

		formattedCmd = append(formattedCmd, part)
	}

	// delete above if we're going with ENVs now

	template.Env = append(template.Env,
		"VOLUME_NAME="+meta.VolumeName,
		"BACKUP_ID="+meta.BackupId,
		"BACKUP_DIR=/backup",
		"DATA_DIR=/data",
	)

	return types.ContainerTemplate{
		Alias:   template.Alias,
		Image:   template.Image,
		Cmd:     formattedCmd,
		Env:     template.Env,
		Volumes: template.Volumes,
	}, nil
}

func waitAndLog(ctx contexts.Context, client client.DockerClient, name string, id string) error {
	waitErr := client.WaitForContainer(name, id)
	logErr := client.ContainerLogs(id)
	return errors.Join(waitErr, logErr)
}

func templateString(template types.ContainerTemplate) string {
	return "alias=" + template.Alias + ",image=" + template.Image + ",cmd='" + strings.Join(template.Cmd, " ") + "',volumes=[" + strings.Join(template.Volumes, "]")
}
