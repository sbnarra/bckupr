# Usage

Bckupr requires at least the backup directory option to run. 

Backups will read [docker labels](labels.md) and can be configured further using the options below.

Restores also requires the backup id but will read labels similar to the backup task and can also be configured further using the options below.

The backup/restore tasks use the same options but with different defaults for the backup id and stop modes, see below for more information.

## Options

Below is a full list of configuration options for the backup and restore tasks.

|Env|Flag|Backup|Restore|
|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|_Required: Must be supplied_|_Required: Must be supplied_|
|`BACKUP_ID`|`--backup-id`|_Optional: Autogenerate Timestamp_|_Required: Must be supplied_|
|`DRY_RUN`|`--dry-run`|_Optional: Defaults `true`_|_Optional: Defaults `true`_|
|`STOP_MODES`|`--stop-modes`|_Optional: Defaults `all,labelled,linked,writers`_|_Optional: Defaults `all,labelled,linked,attached`_|
|`INCLUDE_NAMES`|`--include-names`|_Optional: Defaults none_|_Optional: Defaults none_|
|`INCLUDE_VOLUMES`|`--include-volumes`|_Optional: Defaults none_|_Optional: Defaults none_|
|`EXCLUDE_NAMES`|`--exlclude-names`|_Optional: Defaults none_|_Optional: Defaults none_|
|`EXCLUDE_VOLUMES`|`--exlclude-volumes`|_Optional: Defaults none_|_Optional: Defaults none_|

* `BACKUP_DIR`(_env_)/`--backup-dir`(_cli_)

Sets the location for the backup directory.

* `BACKUP_ID`(_env_)/`--backup-id`(_cli_)

Disable dry runs once finished testing the bckupr deployment.

* `DRY_RUN`(_env_)/`--dry-run`(_cli_)

Disable dry runs once finished testing the bckupr deployment.

### Stop Modes

* `STOP_MODES`(_env_)/`--stop-modes`(_cli_)
    * Valid Values: `all`, `labelled`, `linked`, `writers`, `attached`
    * Backup Default: `labelled`, `linked`, `writers`
    * Restore Default: `labelled`, `linked`, `attached`

Stop modes control how bckupr stops containers before running backups or restores (_backup and restore uses different defaults_).

`all` - Will stop all running containers on the docker host

`labelled` - Will stop containers labelled with `bckupr.stop`, see [Labels](labels.md) for more information.

`linked` - Will stop dependant containers.
* _this option alone has no effect_

`writers` - Will stop all containers with RW access to the volume being backed up or restored.

`attached` - Will stop all containers with the volume being backed up or restored attached.

### Filters

Filters can be applied to container which containers or volumes are picked up during the label scan.

* `INCLUDE_NAMES`(_env_)/`--include-names`(_cli_)
* `INCLUDE_VOLUMES`(_env_)/`--include-volumes`(_cli_)

Include container names and volumes...

* `EXCLUDE_NAMES`(_env_)/`--exclude-names`(_cli_)
* `EXCLUDE_VOLUMES`(_env_)/`--exclude-volumes`(_cli_)

Exclude container names and volumes...