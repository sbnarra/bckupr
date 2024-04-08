# Backup

The `backup` command executes the following process:

1. read container labels to determine volumes to backup, see [labels](../labels.md)
1. stop containers based on [stop modes](#stop_modes)
1. run docker containers to backup each volume into the backup archive within the backup id
1. start containers as each volumes backup completes
1. optionally push completed backup offsite, see [offsite](../offsite.md)

|Env|Flag|Description||
|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|
|`BACKUP_ID`|`--backup-id`|Id of specific backup|_Optional: Autogenerate Timestamp_|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|
|`STOP_MODES`|`--stop-modes`|See [stop modes](#stop_modes) for more info|_Optional: Defaults `labelled,writers,linked`_|

{%
    include-markdown "../../stop-modes.md"
    heading-offset=1
%}

{%
    include-markdown "../../filters.md"
    heading-offset=1
%}
