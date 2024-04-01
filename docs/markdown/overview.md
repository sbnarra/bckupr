# Overview

The bckupr container is a wrapper to the bckupr cli, by default the image executes the `cron` command to run bckupr as a cron server.

By default bckupr runs with `--dry-run=true`/`DRY_RUN=true`, make sure to test and review with dry run enabled before disabling and performing real backups/restores.

The following sections provide an overview of how bckupr works.

## How Backups/Restores Work?

Bckupr is responsible for managing data backups and restores by running customisable containers which mounts the `/backup`(_backup archive_) and `/data`(_volume/mount data_).

The default backup process follows these simple steps...

1. Reading container labels to determine which volumes/mounts need backing up
1. Stopping containers with read/write access to backup volumes/mounts or labelled with `bckupr.stop=true`
    * _this can be configured using the [`--stop-modes`/`STOP_MODES`](stop_modes.md) options._
1. Executing docker containers per volume/mount responsible for performing the backup
    * _this can be configured using `backup` [container templates](container_templates.md)._
1. Finally restarting previously stopped containers.

Bckupr restores follows the same process as the backups but stopping containers with both read/write and read-only access and using `restore` [container templates](container_templates.md).

## Global Options

* `--dry-run`/`DRY_RUN` - _default=`true`_
    * set to `false` to disable dry runs and perform real backups/restores
* `--debug`/`DEBUG` - _default=`false`_
    * set to `true` to enable debug logs
* `--dev`/`DEV` - _default=`false`_
    * set to `true` to enrich logs with dev data (_useful when submitting issues_)
