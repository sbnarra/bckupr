# Labels

Docker labels are used to configure which volumes are backed up alone with additional behaviour.

Bckupr will read the labels of all containers so if a container has crashed or exited, its labels will still be read. If volumes for short lived containers need backing up, you should label the bckupr contanier.

## Options

* `bckupr.volumes=named_vol_1,named_vol_2`

Named volumes can configured using the above label with a comma separated list.

* `bckupr.volumes.<alias>=/path/to/volume/mount`

Mounted volumes require an alias so should be configured using the label above per-mount.

* `bckupr.stop=true|false`

By default bckupr will stop all containers, and their dependancies, with write access using `stop modes`. If this behaviour is change, the above label can be used to tell bckupr to always shutdown the container.

## Prefix

* `LABEL_PREFIX`(_env_)/`--label-prefix`(_cli_)

The label prefix `bckupr` can be changed using the above option, this may be useful if running multiple instances.

Running bckupr with `--label-prefix=custom` will result it in scanning for labels like `custom.volumes=`.