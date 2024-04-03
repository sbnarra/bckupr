# Notifications

Notifications can be enabled using [Shoutrrr](https://containrrr.dev/shoutrrr).

## Configuration

* `NOTIFICATION_URLS`(_env_)/`--notification-urls`(_cli_)

The above option is required to configure notifications, see the [shoutrrr documentation](https://containrrr.dev/shoutrrr/latest/services/overview/) on how to configure your services.

### Tuning

By default bckupr will send notifications for all backup/restore jobs that are started/completed and notifications per-volume. The below flags can be used to tune which notifications you recieve.

---

Options for configuring backup/restore started/completed notifications:

* `NOTIFY_JOB_STARTED`(_env_)/`--notify-job-started`(_cli_)
* `NOTIFY_JOB_COMPLETED`(_env_)/`--notify-job-completed`(_cli_)
* `NOTIFY_JOB_ERROR`(_env_)/`--notify-job-error`(_cli_)

Options for configuring per-volume started/completed notifications:

* `NOTIFY_TASK_STARTED`(_env_)/`--notify-task-started`(_cli_)
* `NOTIFY_TASK_COMPLETED`(_env_)/`--notify-task-completed`(_cli_)
* `NOTIFY_TASK_ERROR`(_env_)/`--notify-task-error`(_cli_)