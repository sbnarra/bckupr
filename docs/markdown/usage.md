# Usage

Bckupr is designed to run as a daemon process using the [`daemon`](command/daemon.md) command.

Once started you can then run the [`backup`](command/backup.md), [`restore`](command/restore.md), [`list`](command/list.md), [`delete`](command/delete.md) and [`rotate`](command/rotate.md) cli commands or use the UI to execute the same functionality.

* _The daemon isn't required, the client commands can be run without the daemon using the `--no-daemon` flag._