package discover

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

var UnableToDetect = errors.New("Unable to Detect Host Backup Dir")

func MountedBackupDir(ctx contexts.Context, dockerHosts []string) (string, *errors.Error) {
	var fErr *errors.Error
	for _, dockerHost := range dockerHosts {
		dir, err := mountedBackupDir(ctx, dockerHost)
		if errors.Is(err, UnableToDetect) {
			logging.CheckWarn(ctx, err, dockerHost)
		} else {
			fErr = errors.Wrap(err, dockerHost)
		}

		if dir != "" {
			return dir, nil
		}
	}
	return "", errors.Wrap(fErr, "supply --"+keys.HostBackupDir.CliId)
}

func mountedBackupDir(ctx contexts.Context, dockerHost string) (string, *errors.Error) {
	if val := os.Getenv("BCKUPR_IN_CONTAINER"); val != "1" {
		return "", errors.Wrap(UnableToDetect, "not running in container")
	}
	version := os.Getenv("VERSION")

	docker, err := client.Client(ctx, dockerHost)
	if err != nil {
		return "", err
	}

	var c *dockerTypes.Container

	kv := func(key, value string) filters.KeyValuePair {
		return filters.KeyValuePair{Key: key, Value: value}
	}
	if found, err := docker.FindContainers(ctx,
		kv("label", "org.opencontainers.image.ref.name=sbnarra/bckupr"),
		kv("label", "org.opencontainers.image.version="+version),
		kv("volume", ctx.ContainerBackupDir),
	); err != nil {
		return "", err
	} else if foundLen := len(found); foundLen == 1 {
		c = &found[0]
	} else if foundLen > 1 {
		if c, err = detectRunningInstance(ctx, docker, found); err != nil {
			return "", err
		}
	} else {
		return "", errors.Wrap(UnableToDetect, "bckupr container not matched with labels")
	}

	if c == nil {
		return "", errors.Wrap(UnableToDetect, "bckupr container not found")
	}

	backupDirHostDir := backupDirHostDir(c, ctx.ContainerBackupDir)
	return backupDirHostDir, nil
}

func detectRunningInstance(ctx contexts.Context, docker client.DockerClient, cs []dockerTypes.Container) (*dockerTypes.Container, *errors.Error) {
	detectionFile := detectionFile()

	for _, c := range cs {
		if err := docker.Exec(ctx, c.ID, []string{"touch", detectionFile}, true); err != nil {
			logging.CheckError(ctx, err)
			continue
		}

		_, err := os.Stat(detectionFile)
		go docker.Exec(ctx, c.ID, []string{"rm", detectionFile}, true)
		if err != nil {
			continue
		} else {
			return &c, nil
		}

	}
	return nil, nil
}

func backupDirHostDir(c *dockerTypes.Container, containerBackupDir string) string {
	for _, mount := range c.Mounts {
		if mount.Destination == containerBackupDir {
			return mount.Source
		}
	}
	return ""
}

func detectionFile() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("/tmp/%X", b)
}
