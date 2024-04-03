# Offsite

...incomplete...

Bckupr can push backups offsite to any storage option using custom docker containers.

These can be defined with YAML and mounted to bckupr. Bckupr also comes with predefined config for AWS, SCP, Azure, GCP.

* ___Future Work: Add templates for SCP, Azure, GCP___

The option below should be set to the path of the offsite container config.

* `OFFSITE_CONTAINERS`(_env_)/`--offsite-containers`(_cli_)

_Predefined configs can be used:_

* `/offsite/scp.yml`
* `/offsite/aws-s3.yml`
* `/offsite/azure-.yml`
* `/offsite/gcp-.yml`

___Future Work: Explain how to define custom templates___