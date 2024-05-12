<h1 align="center">
  Bckupr
</h1>

<p align="center">
Bckupr is program to automate backup creation and data restoration.
</p>
<p align="center">
This tool reads docker labels to determine which volumes/mounts require backing up, cleanly stopping containers and their dependancies to make sure all data is flushed to disk to create consistent backups.
</p>

## Quick Start

Using Bckupr you can automate local backups, pushing to offsite storage and data rentention with simple commands to also automate restoring your data. 

Bckupr will read container labels to tell which volumes should be backed up before shutting down relavent containers and performing backups to ensure all data is flushed to disk avoiding corrupt backups.

<!-- https://docs.docker.com/storage/volumes/#back-up-restore-or-migrate-data-volumes -->

To get started simple tag your containers with the volume to backup using `bckupr.volumes=<volume-name>`.

Next run Bckupr using the following docker commands:

=== "docker run"
    Use the following docker run command to start bckupr:
    ```bash
    $ docker run --name bckupr -d \
        -p 8000:8000 \
        -v /tmp/backups:/backups \
        -v /var/run/docker.sock:/var/run/docker.sock \
        sbnarra/bckupr
    ```
    Create new ad-hoc backup:
    ```bash
    $ docker exec bckupr backup
    ```
    To then restore from the adhoc backup:
    ```bash
    $ docker exec bckupr restore --backup-id <id-from-backup-logs>
    ```
=== "docker-compose.yml"
    Use the following YAML to run bckupr:
    ```yaml
    version: "3"
    services:
      bckupr:
        image: sbnarra/bckupr
        ports:
          - 8000:8000
        volumes:
          - /tmp/backups:/backups
          - /var/run/docker.sock:/var/run/docker.sock
    ```
    Create new ad-hoc backup:
    ```bash
    $ docker compose bckupr exec backup
    ```
    To then restore from the adhoc backup:
    ```bash
    $ docker compose bckupr exec restore --backup-id <id-from-backup-logs>
    ```

_By default bckupr runs in dry run mode, to disable use arg `--dry-run false` or env `DRY_RUN=false` once testing is complete._

_Don't forget to update `/tmp/backups` to backup archieve._