# Notifications

Notifications can be enabled using [Shoutrrr](https://containrrr.dev/shoutrrr).

## Configuration

The above option is required to configure notifications, see the [shoutrrr documentation](https://containrrr.dev/shoutrrr/latest/services/overview/) on how to configure your services.

|Env|Flag|Description|Default|
|-|-|-|-|
|`NOTIFICATION_URLS`|`--notification-urls`|Comma separated list of [service urls](https://containrrr.dev/shoutrrr/latest/services/overview/)|_Optional: Defaults none_|

### Tuning

By default bckupr will send notifications for all backup/restore jobs that are started/completed and notifications per-volume. The flags below can be used to tune which notifications you recieve.

|Env|Flag|Description|Default|
|-|-|-|-|
|`NOTIFY_JOB_STARTED`|`--notify-job-started`|Enables debug logging|_Optional: Defaults `true`_|
|`NOTIFY_JOB_COMPLETED`|`--notify-job-completed`|Enables debug logging|_Optional: Defaults `true`_|
|`NOTIFY_JOB_ERROR`|`--notify-job-error`|Enables debug logging|_Optional: Defaults `true`_|
|`NOTIFY_TASK_STARTED`|`--notify-task-started`|Enables debug logging|_Optional: Defaults `true`_|
|`NOTIFY_TASK_COMPLETED`|`--notify-task-completed`|Enables debug logging|_Optional: Defaults `true`_|
|`NOTIFY_TASK_ERROR`|`--notify-task-error`|Enables debug logging|_Optional: Defaults `true`_|
