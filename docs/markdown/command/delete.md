# Delete

The `delete` command permanently removes a local backup from disk (_this doesn't remove offsite backups_).

|Env|Flag|Description||
|-|-|-|-|
|`BACKUP_ID`|`--backup-id`|Id of specific backup|_Required: Must be supplied_|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|