# Stop Modes

Stop modes control how bckupr chooses which containers needs shutting down, multiple can be applied at once. There's 4 stop modes:

* `all` - stop all running containers
* `labelled` - stop all containers labelled with `bckupr.stop`
* `writers` - stop all containers with backup/restore volumes mounted with read/write access
* `attached` - stop all containers with backup/restore volumes mounted

By default backups use `labelled` and `writers` stop modes. Restores use `labelled` and `attached` stop modes.
