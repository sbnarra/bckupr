# Cron

The `cron` command automates the `backup` and `rotate` commands using cron expressions.

## Backup

By default the backup schedule is enabled with `0 0 * * *`, daily backups at midnight.

|Env|Flag|Description|Default|
|-|-|-|-|
|`BACKUP_SCHEDULE`|`--backup-schedule`|Sets schedule for automatic backup creation, set to "" to disable|_Optional: Defaults `0 0 * * *`_|

See the [`backup`](backup.md) command on how to configure backups.

## Rotate

By default the rotate schedule is disabled.

|Env|Flag|Description|Default|
|-|-|-|-|
|`ROTATE_SCHEDULE`|`--rotate-schedule`|Sets schedule for automatic backup rotation, set to "" to disable|_Optional: Defaults none_|

See the [`rotate`](rotate.md) command on how to configure backup rotations.