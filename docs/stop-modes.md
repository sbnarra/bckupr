# Stop Modes

Stop modes control how bckupr stops containers before running backups or restores.

`all` - Will stop all running containers on the docker host

`labelled` - Will stop containers labelled with `bckupr.stop`, see [Labels](./markdown/labels.md) for more info.

`writers` - Will stop all containers with RW access to the volume.

`attached` - Will stop all containers with the volume attached.

`linked` - Will stop dependant containers of those targeted from `labelled`, `wrtiers` or `attached` (_this option alone has no effect_).
