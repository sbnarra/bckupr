# Offsite

Bckupr can push backups offsite to any storage option using custom docker containers and then pull those backups if deleted locally.

These can be defined with YAML and mounted into the bckupr container. Bckupr also comes with predefined config:
<!-- * SCP -->
* AWS - `/offsite/aws-s3.yml`
<!-- * Azure -->
<!-- * GCP -->

Each integration will require environment variables sets in the bckupr container which will be passed into the offsite container during backup pushes or restore pulls.

# Configure

|Env|Flag|Description|Default|
|-|-|-|-|
|`OFFSITE_CONTAINERS`|`--offsite-containers`|Path to config file|_Optional: Defaults none_|

## Predefined

### SCP

* `OFFSITE_CONFIG=/offsite/scp.yml`

Requires the environment variables:

* `HOSTNAME`
* `SSH_DIR`

### AWS S3

* `OFFSITE_CONFIG=/offsite/aws-s3.yml`

Requires the environment variables:

* `SECRET_KEY`
* `REGION`
* `BUCKET`

<!-- ### Azure

* `OFFSITE_CONFIG=/offsite/azure-.yml`

Requires the environment variables:

* 

### GCP

* `OFFSITE_CONFIG=/offsite/gcp-.yml`

Requires the environment variables:

*  -->


## Custom

Configuration to push/pull from any offsite storage can be defined using the template below:

```yaml
push:
  image: ubuntu
  cmd:
    - sh
    - -c
    - scp $BACKUP_PATH $USERNAME@$HOSTNAME:/backups/$BACKUP_ID/$VOLUME_NAME$EXT
  env:
    - USERNAME
    - HOSTNAME
pull:
  image: ubuntu
  cmd:
    - sh
    - -c
    - scp $USERNAME@$HOSTNAME:/backups $BACKUP_PATH
  env:
    - USERNAME
    - HOSTNAME
```

`image`: image to be used.

`cmd`: the command is made up of a list of arguments.

`env`: this is a list of environment variable names to be passed from the bckupr container into the offsite containers.

Mount the file into the bckupr container and pass the flag `--offsite-containers=/path/to/file.yml`