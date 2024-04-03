package containers

import (
	"errors"
	"strings"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func (c *Containers) RunContainer(ctx contexts.Context, meta RunMeta, template types.ContainerTemplate) error {
	logging.Debug(ctx, "Template:", templateString(template))
	if rendered, err := renderTemplate(template, meta); err != nil {
		return err
	} else {
		return c.runContainer(ctx, rendered)
	}
}

func (c *Containers) runContainer(ctx contexts.Context, rendered types.ContainerTemplate) error {
	if ctx.DryRun {
		logging.Info(ctx, "Dry Run!", templateString(rendered))
		return nil
	}
	logging.Info(ctx, "Executing:", templateString(rendered))

	id, err := c.client.RunContainer(ctx, rendered.Image, rendered.Cmd, rendered.Env, rendered.Volumes)
	if err == nil {
		err = c.waitAndLog(ctx, rendered.Alias, id)
	}

	removalErr := c.client.RemoveContainer(ctx, id)
	return errors.Join(err, removalErr)
}

func renderTemplate(template types.ContainerTemplate, meta RunMeta) (types.ContainerTemplate, error) {
	if len(template.Image) == 0 || len(template.Cmd) == 0 {
		return types.ContainerTemplate{}, &MissingTemplate{Message: "No config for " + template.Alias}
	}

	template.Env = append(template.Env,
		"VOLUME_NAME="+meta.VolumeName,
		"BACKUP_ID="+meta.BackupId,
		"BACKUP_DIR=/backup",
		"DATA_DIR=/data",
	)

	return types.ContainerTemplate{
		Alias:   template.Alias,
		Image:   template.Image,
		Cmd:     template.Cmd,
		Env:     template.Env,
		Volumes: template.Volumes,
	}, nil
}

func (c *Containers) waitAndLog(ctx contexts.Context, name string, id string) error {
	waitErr := c.client.WaitForContainer(ctx, name, id)
	logErr := c.client.ContainerLogs(ctx, id)
	return errors.Join(waitErr, logErr)
}

func templateString(template types.ContainerTemplate) string {
	return "alias=" + template.Alias + ",image=" + template.Image + ",cmd='" + strings.Join(template.Cmd, " ") + "',volumes=[" + strings.Join(template.Volumes, "]")
}
