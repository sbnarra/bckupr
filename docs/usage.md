# Usage

Bckupr comes as a single binary.

To run bckupr as a long running process there's the `daemon` command which runs the bckupr service/cron/gui, and `cron` which only runs scheduled backups with no gui.

Once running the bckupr `daemon` process you can then use the `backup`, `restore`, `list`, `rotate` and `delete` commands.

_The commands can be run without first running the daemon process by using the `--no-daemon` flag_ 

---

Bckupr requires at least the backup directory option to run. 

Backups will read [docker labels](labels.md) and can be configured further using the options below.

Restores also requires the backup id but will read labels similar to the backup task and can also be configured further using the options below.

The backup/restore tasks use the same options but with different defaults for the backup id and stop modes, see below for more information on the different options for different tasks.

_Global_:

|Env|Flag|Description|Backup|Restore|
|-|-|-|-|-|
|`DEBUG`|`--debug`|Enables debug logging|_Optional: Defaults `false`_|_Optional: Defaults `false`_|

## Daemon



|Env|Flag|Description|||
|-|-|-|-|-|
|`UNIX_SOCKET`|`--unix-socket`|Enables debug logging|_Optional: Defaults `false`_|
|`TCP_ADDR`|`--tcp-addr`|Enables debug logging|_Optional: Defaults `false`_|
|`EXPOSE_API`|`--expose-api`|Enables debug logging|_Optional: Defaults `false`_|
|`UI_ENABLED`|`--ui-enabled`|Enables debug logging|_Optional: Defaults `false`_|
|`METRICS_ENABLED`|`--metrics-enabled`|Enables debug logging|_Optional: Defaults `false`_|

_`backup`, `restore` and `rotate` options are supported to set default values._

## Cron


|Env|Flag|Description|||
|-|-|-|-|-|
|`TIMEZONE`|`--unix-socket`|Enables debug logging|_Optional: Defaults `false`_|
|`BACKUP_SCHEDULE`|`--unix-socket`|Enables debug logging|_Optional: Defaults `false`_|
|`ROTATE_SCHEDULE`|`--unix-socket`|Enables debug logging|_Optional: Defaults `false`_|

See the [Cron docs](cron.md) for more information.

## Backup/Restore

These are the main options for backup/restore tasks.

|Env|Flag|Description|Backup|Restore|
|-|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|_Required: Must be supplied_|
|`BACKUP_ID`|`--backup-id`|Id of specific backup|_Optional: Autogenerate Timestamp_|_Required: Must be supplied_|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|_Optional: Defaults `true`_|

### Stop Modes

Stop modes control how bckupr stops containers before running backups or restores (_backup and restore uses different defaults_).

`all` - Will stop all running containers on the docker host

`labelled` - Will stop containers labelled with `bckupr.stop`, see [Labels](labels.md) for info on labels.

`writers` - Will stop all containers with RW access to the volume being backed up or restored.

`attached` - Will stop all containers with the volume being backed up or restored attached.

`linked` - Will stop dependant containers (_this option alone has no effect_).

|Env|Flag|Backup|Restore|
|-|-|-|-|
|`STOP_MODES`|`--stop-modes`|_Optional: Defaults `labelled,writers,linked`_|_Optional: Defaults `labelled,attached,linked`_|


### Filters

Filters can be applied to limit which containers/volumes are included with the backup/restore tasks.

|Env|Flag|Backup/Restore|
|-|-|-|
|`INCLUDE_NAMES`|`--include-names`|_Optional: Defaults none_|
|`INCLUDE_VOLUMES`|`--include-volumes`|_Optional: Defaults none_|
|`EXCLUDE_NAMES`|`--exlclude-names`|_Optional: Defaults none_|
|`EXCLUDE_VOLUMES`|`--exlclude-volumes`|_Optional: Defaults none_|

## List

List backups with metadata.

|Env|Flag|Description||
|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|_Required: Must be supplied_|

## Delete

Delete existing backup.

|Env|Flag|Description||
|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|
|`BACKUP_ID`|`--backup-id`|Id of specific backup|_Required: Must be supplied_|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|