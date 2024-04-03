# Offsite

Bckupr can push backups offsite to any storage option using custom docker containers and then pull those backups if deleted locally.

These can be defined with YAML and mounted into the bckupr container. Bckupr also comes with predefined config for AWS, SCP, Azure, GCP.

Each integration will require environment variables sets in the bckupr container which will be passed into the offsite container during backup pushes or restore pulls.

* ___Future Work: Add templates for SCP, Azure, GCP___

# Option

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

* 

### Azure

* `OFFSITE_CONFIG=/offsite/azure-.yml`

Requires the environment variables:

* 

### GCP

* `OFFSITE_CONFIG=/offsite/gcp-.yml`

Requires the environment variables:

* 


## Custom

___Future Work: Explain how to define custom templates___

```yaml
push:
    image: <container-image>
    cmd: 
    env:
```