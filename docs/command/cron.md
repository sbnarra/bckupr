# Cron

Bckupr uses a builtin cron scheduler for automation of backup creation and backup rotation.

## Backup Creation

By default the backup schedule is enabled with `0 0 * * *`, daily backups at midnight.

The scheduled backups will inherit the options supplied at startup, see the [usage docs](usage.md) for more information.

|Env|Flag|Description|Default|
|-|-|-|-|
|`BACKUP_SCHEDULE`|`--backup-schedule`|Sets schedule for automatic backup creation, set to "" to disable|_Optional: Defaults `0 0 * * *`_|

## Backup Rotation

|Env|Flag|Description|Default|
|-|-|-|-|
|`ROTATE_SCHEDULE`|`--rotate-schedule`|Sets schedule for automatic backup rotation, set to "" to disable|_Optional: Defaults none_|