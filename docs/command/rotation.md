# Rotation

The `rotate` command is used to clean up old backups. You configure policies which define how many backups to retain for periods of time. 

## Modes

By default `rotate` runs in the dry run mode, it's highly recommended to test your policies throughly using dry runs. Once testing is finished use the `DRY_RUN=false`/`--dry-run=false` option.

Once dry runs is disabled, `rotate` won't delete backups but move them into a _bin_ directory within the backups directory, this directly needs emptying manually. To skip the `bin` directory and delete backups you'll need to supply the `DESTROY_BACKUPS=true`/`--destroy-backups=true` option.

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
      from: -7d # 7 days
      to: -1d
    keep: 1
  # final catch all
  - period:
      from: -999w
      to: -7d # 7 days
    keep: -1
```

This configuration will:
1. keep the most recent backup in the last 24 hours
  * _all other backups within the 24hour window are removed_
1. keep the most recent backup in the last 6 days
  * _these 2 policies covers the last week, specifify daily limit and weekly limit. Notice the 2nd policy doesn't overlap with the 1st_
1. final catch all policy will keep the oldest backup using a negative keep
  * _careful keep oldest as you'll be removing newest (should move this example in the 1st policy)_