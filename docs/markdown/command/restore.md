# Restore

The `restore` command requires a valid `backup-id` and executes the following process:

1. check backup exists locally
    * if not found, optionally pull backup from offsite, see [offsite](../offsite.md)
1. stop containers based on [stop modes](#stop_modes)
1. run docker containers to restore each volume from the backup archive within the supplied backup id
1. start containers as each volumes restore completes

|Env|Flag|Description||
|-|-|-|-|
|`BACKUP_DIR`|`--backup-dir`|Directory containing local backups|_Required: Must be supplied_|
|`BACKUP_ID`|`--backup-id`|Id of specific backup|_Required: Must be supplied_|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|
|`STOP_MODES`|`--stop-modes`|See [stop modes](#stop_modes) for more info|_Optional: Defaults `labelled,attached,linked`_|

{%
    include-markdown "../../stop-modes.md"
    heading-offset=1
%}

{%
    include-markdown "../../filters.md"
    heading-offset=1
%}
