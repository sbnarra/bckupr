# Backup Options

Backup configuration explained...

* `BACKUP_DIR`(_env_)/`--backup-dir`(_cli_)
    * _required_

...

* `DRY_RUN`(_env_)/`--dry-run`(_cli_)
    * Valid Values: `true`

...

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

## Filters

Filters can be applied to container which containers or volumes are picked up during the label scan.

* `INCLUDE_NAMES`(_env_)/`--include-names`(_cli_)
* `INCLUDE_VOLUMES`(_env_)/`--include-volumes`(_cli_)

Include container names and volumes...

* `EXCLUDE_NAMES`(_env_)/`--exclude-names`(_cli_)
* `EXCLUDE_VOLUMES`(_env_)/`--exclude-volumes`(_cli_)

Exclude container names and volumes...