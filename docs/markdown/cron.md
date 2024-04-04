# Cron

Bckupr uses a builtin cron scheduler to automate your backups. By default this schedule is set to `0 0 * * *`, daily backups at midnight.

The scheduled backups will use inherit the options supplied at startup, see the [usage docs](usage.md) for more information.

|Env|Flag|Description|Default|
|-|-|-|-|
|`BACKUP_SCHEDULE`|`--backup-schedule`|Sets schedule for automatic backups, set to "" to disable|_Optional: Defaults `0 0 * * *`_|
