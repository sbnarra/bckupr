# Cron

Bckupr can backup your volumes using a cron schedule configurable using the `BACKUP_SCHEDULE`(_env_)/`--backup-schedule`(_cli_) options.

By default this schedule is set to `0 0 * * *`, daily backups at midnight. Set this schedule to "" to disable.

The scheduled backups will use inherit the options supplied at startup, see the [backup options](options.md) for more information.