# Daemon

The `daemon` command is a long running task that should be used to start bckupr as a background service. You can then use the [`backup`](backup.md), [`restore`](restore.md), [`rotate`](rotate.md), [`list`](list.md), and [`delete`](delete.md) commands or use the UI on http://localhost:8000.

This also runs the builtin cron scheduler to automate executing the `backup` and `rotate` commands. To run bckupr as only a cron instance see the [`cron`](cron.md) command.

|Env|Flag|Description|||
|-|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|
|`UNIX_SOCKET`|`--unix-socket`|Path to bckupr unix socket, used by cli commands|_Optional: Defaults `.bckupr.sock`_|
|`TCP_ADDR`|`--tcp-addr`|Tcp bind address|_Optional: Defaults `0.0.0.0:8000`_|
|`EXPOSE_API`|`--expose-api`|Exposes API via TCP for external access|_Optional: Defaults `false`_|
|`UI_ENABLED`|`--ui-enabled`|Enables bckupr GUI|_Optional: Defaults `true`_|
|`METRICS_ENABLED`|`--metrics-enabled`|Enables metrics, see [metrics](../metrics.md) for more info|_Optional: Defaults `false`_|

This command also accepts the [`cron`](cron.md), [`backup`](backup.md) and [`rotate`](rotate.md) options.