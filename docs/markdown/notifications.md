# Notifications

Notifications can be enabled via bckuprs [Shoutrrr](https://containrrr.dev/shoutrrr) intergration.

## Configuration

Set the option below to configure notifications, see the [shoutrrr documentation](https://containrrr.dev/shoutrrr/latest/services/overview/) on how to configure your services.

|Env|Flag|Description|Default|
|-|-|-|-|
|`NOTIFICATION_URLS`|`--notification-urls`|Comma separated list of [service urls](https://containrrr.dev/shoutrrr/latest/services/overview/)|_Optional: Defaults none_|

### Tuning

By default bckupr will send notifications for all backup/restore jobs that are started/completed and notifications per-volume. The flags below can be used to tune which notifications you recieve.

|Env|Flag|Description|Default|
|-|-|-|-|
|`NOTIFY_JOB_STARTED`|`--notify-job-started`|Notify when backup/restore starts|_Optional: Defaults `true`_|
|`NOTIFY_JOB_COMPLETED`|`--notify-job-completed`|Notify when backup/restore completes(_success & error_)|_Optional: Defaults `true`_|
|`NOTIFY_JOB_ERROR`|`--notify-job-error`|Notify when backup/restore errors|_Optional: Defaults `true`_|
|`NOTIFY_TASK_STARTED`|`--notify-task-started`|Notify when volume backup/restore starts|_Optional: Defaults `true`_|
|`NOTIFY_TASK_COMPLETED`|`--notify-task-completed`|Notify when volume backup/restore completes(_success & error_)|_Optional: Defaults `true`_|
|`NOTIFY_TASK_ERROR`|`--notify-task-error`|Notify when volume backup/restore errors|_Optional: Defaults `true`_|
