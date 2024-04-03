# Cron

Bckupr uses a builtin cron scheduler to automate your backups. By default this schedule is set to `0 0 * * *`, daily backups at midnight.

The scheduled backups will use inherit the options supplied at startup, see the [backup options](options.md) for more information.

## Options

* `BACKUP_SCHEDULE`(_env_)/`--backup-schedule`(_cli_)

    _default_: ``0 0 * * *``

Override the backup scheduler or set this to "" to disable.
