# Bckupr

Backup automation and management for docker volumes using labelled containers.

```bash
$ docker run --name bckupr -d \
    -p 8000:8000 \
    -e BACKUP_DIR=/tmp/backups \
    -v /tmp/backups:/tmp/backups \
    -v /var/run/docker.sock:/var/run/docker.sock \
    sbnarra/bckupr
```

Label containers with `bckupr.volumes=<named_volume>` or `bckupr.volumes.<mount_alias>=<mount_path>`

Access UI at http://localhost:8000

See [User Documentation](https://sbnarra.github.io/bckupr) for more information.