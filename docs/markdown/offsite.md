# Offsite

Bckupr can push backups offsite to any storage option using custom docker containers.

These can be defined with YAML and mounted to bckupr. Bckupr also comes with predefined config for AWS, _??? SCP, Azure, GCP ???_.

* `OFFSITE_CONTAINERS`(_env_)/`--offsite-containers`(_cli_)

The above option should be set to the path of the offsite container config.

Predefined config can be used:
* `/offsite/aws-s3.yml`
* `/offsite/scp.yml`
* `/offsite/.yml`
* `/offsite/aws-s3.yml`

```yaml
push:
  image: amazon/aws-cli
  env: 
    - AWS_SECRET
    - AWS_SECRET_KEY
  cmd:
    - aws
    - s3
    - cp
    - /backup/{backup_id}/{name}.tar.gz
    - s3://$BUCKET:{backup_id}/{name}.tar.gz
pull:
  image: amazon/aws-cli
  env: 
    - AWS_SECRET
    - AWS_SECRET_KEY
  cmd:
    - aws
    - s3
    - cp
    - s3://$BUCKET:{backup_id}/{name}.tar.gz
    - /backup/{backup_id}/{name}.tar.gz
```