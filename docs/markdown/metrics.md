# Metrics

___Future Work: Use backup/restore specific metrics___

Bckupr can export metrics via the `/metrics` endpoint to be scrapped by prometheus.

|Env|Flag|Description|Default|
|-|-|-|-|
|`METRICS_ENABLED`|`--metrics-enabled`|Enables `/metrics` endpoint|_Optional: Defaults `false`_|

## Backup

* `backup_duration_seconds`
* `backup_success_total`
    * `id`: backup id
    * `volume`: volume path
* `backup_error_total`
    * `id`: backup id
    * `volume`: volume path
    * `error`: error message

## Restore

* `restore_duration_seconds`
* `restore_success_total`
    * `id`: backup id
    * `volume`: volume path
* `restore_error_total`
    * `id`: backup id
    * `volume`: volume path
    * `error`: error message

