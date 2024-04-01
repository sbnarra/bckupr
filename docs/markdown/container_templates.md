# Container Templates

Bckupr itself doesn't perform any backups or restores, it insteads runs containers based on templates to do the job for it.

By default there's 3 containers pre-configured for created/restoring compressed tar backups, [`file-check`](#file-check), [`backup`](#backup) and [`restore`](#restore); this follows https://docs.docker.com/desktop/backup-and-restore/.

Then there's 2 unconfigured templates for handling offsite backups:

* `offsite-push`: triggered after the `backup` container completes
* `offsite-pull`: triggers if `file-check` fails during the restore process

Templates accept 3 options; image, cmd and env.

* `--<name>-image`/`<NAME>_IMAGE`: e.g `--backup-image`/`BACKUP_IMAGE` - sets containers image
* `<name>-cmd`: e.g `--backup-cmd`/`BACKUP_CMD` - sets command to execute within container
* `<name>-env`: e.g `--backup-env`/`BACKUP_ENV` - sets env vars to passthrough

## Template Replacements

* `{name}`
* `{backup_id}`

## Pre-Configured Templates

### backup

Responsible for creating a _local_ backup.

* `--backup-image`/`BACKUP_IMAGE` = `busybox`
* `--backup-cmd`/`BACKUP_CMD` = `tar czvf {name}.tar.gz -C /backup/{backup_id} /data`
* `--backup-env`/`BACKUP_ENV` = _N/A_

### file-check

Responsible for checking `backup` created the expected backup file and the backup exists before `restore` executes.

* `--file-check-image`/`BACKUP_IMAGE` = `busybox`
* `--file-check-cmd`/`BACKUP_CMD` = `ls /backup/{backup_id}/{name}.tar.gz`
* `--file-check-env`/`BACKUP_ENV` = _N/A_

### restore

Responsible for restoring data from a _local_ backup.

* `--restore-image`/`RESTORE_IMAGE` = `busybox`
* `--restore-cmd`/`RESTORE_CMD` = `tar xzvf /backup/{backup_id}/{name}.tar.gz --strip 1 -C /data`
* `--restore-env`/`RESTORE_ENV` = _N/A_

## Un-Configured Templatess

### offsite-push

### offsite-pull