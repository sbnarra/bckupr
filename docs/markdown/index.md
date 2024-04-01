<h1 align="center">
  Bckupr
</h1>

<p align="center">
The all in one container backup and restore solution.
</p>

## Quick Start

Using Bckupr you can automate local backups, pushing to offsite storage and data rentention with simple commands to also automate restoring your data. Bckupr will read container labels to tell which volumes should be backed up before shutting down relavent containers and taking backups to ensure all data is flushed to disk avoiding corrupt backuprs. Run Bckupr using the following docker commands to get started...

=== "docker run"
    Use the following docker run command to start bckupr:
    ```bash
    $ docker run --name bckupr -d \
        -v /var/run/docker.sock:/var/run/docker.sock \
        sbnarra/bckupr cron \
        --backup-host-dir /backups \
        --schedule "0 0 * * *"
    ```
    Create new ad-hoc backup:
    ```bash
    $ docker exec bckupr backup --backup-id adhoc
    ```
    To then restore from the adhoc backup:
    ```bash
    $ docker exec bckupr restore --backup-id adhoc
    ```
=== "docker-compose.yml"
    Use the following YAML to run bckupr:
    ```yaml
    version: "3"
    services:
      bckupr:
        image: sbnarra/bckupr
        environment:
          BACKUP_HOST_DIR: /backups
          SCHEDULE: "0 0 * * *"
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock
    ```
    Create new ad-hoc backup:
    ```bash
    $ docker compose bckupr exec backup --backup-id adhoc
    ```
    To then restore from the adhoc backup:
    ```bash
    $ docker compose bckupr exec restore --backup-id adhoc
    ```
_By default bckupr runs in dry run mode, to disable use arg `--dry-run false` or env `DRY_RUN=false` once testing is complete._