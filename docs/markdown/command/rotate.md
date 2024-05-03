# Rotate

The `rotate` command is used to clean up old backups. You configure policies which define how many backups to retain for a period of time. These are then moved into a bin directory or simply deleted from disk (_has no affect on offsite backups_).

|Env|Flag|Description||
|-|-|-|-|
|`DRY_RUN`|`--dry-run`|Needs disabling once completed testing|_Optional: Defaults `true`_|
|`DESTROY_BACKUPS`|`--destroy-backups`|Backups are removed from disk rather than moved to _bin_|_Optional: Defaults `true`_|

## Modes

By default `rotate` runs in the dry run mode, it's highly recommended to test your policies throughly using dry runs. Once testing is finished use the `DRY_RUN=false`/`--dry-run=false` option.

Once dry runs are disabled, `rotate` won't delete backups but move them into a _bin_ directory within the backups archive directory, this directory needs emptying manually. To skip the `bin` directory and automatically delete backups you need to supply the `DESTROY_BACKUPS=true`/`--destroy-backups=true` option.

## Policies

Example policies configuration:
```yaml
policies:
  # policy for last 24 hours
  - period:
      from: -1d
      to: 0s
    keep: 1
  # policy for last week
  - period:
      from: -7d
      to: -1d
    keep: 1
  # final catch all
  - period:
      from: -999w
      to: -7d
    keep: -1
```

This configuration will:

1. keep the most recent backup in the last 24 hours
    * _all other backups created within the 24 hour window are removed_
1. keep the most recent backup in the last 6 days
  * _these 2 policies covers the last week, specifify daily limit and weekly limit. Notice the 2nd policy doesn't overlap with the 1st_
1. final catch all policy will keep the oldest backup using a negative keep
  * _careful keep oldest as you'll be removing newest (should move this example in the 1st policy)_