check_data_is() {
    _check_data_is "volume" $@
    _check_data_is "mount" $@
    echo "data matches: '$*'"
}

_check_data_is() {
    vol=$1; shift
    MSG=$(docker exec data_writer cat /mnt/${vol}/msg)
    if [ "$MSG" != "$*" ]; then
        docker logs bckupr
        echo "invalid mount data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
}

write_data() {
    docker exec data_writer sh -c "echo $1 >/mnt/volume/msg"
    docker exec data_writer sh -c "echo $1 >/mnt/mount/msg"
    echo "Wrote data '$1'"
}

### pre test setup ###
mkdir -p $PWD/.test_filesystem/backups
rm -rf /tmp/bckupr/mount
docker rm -f bckupr
docker rm -f data_writer
docker volume rm test_volume_backup

set -e ### starting test ###

mkdir -p /tmp/bckupr/mount
docker volume create test_volume_backup

VERSION=test ./scripts/app-build-image.sh
docker run --name bckupr -d \
    -e DRY_RUN=true -e DEBUG=true \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $PWD/.test_filesystem/backups:/backups \
    sbnarra/bckupr:test

docker run --name data_writer -d \
    -l bckupr.stop=true \
    \
    -l bckupr.volumes=test_volume_backup \
    -v test_volume_backup:/mnt/volume \
    \
    -l bckupr.volumes.test_mount_backup=$PWD/.test_filesystem/example-mount \
    -v $PWD/.test_filesystem/example-mount:/mnt/mount \
    alpine sleep 120

write_data pre-backup

BACKUP_ID="test-$(date +%Y-%m-%d_%H-%M)"
docker exec bckupr bckupr backup --dry-run=false --backup-id=$BACKUP_ID

write_data post-backup

docker exec bckupr bckupr restore --dry-run=false --include-names data_writer --backup-id $BACKUP_ID

check_data_is pre-backup

### post test cleanup ###
echo "test completed successfully, running clean up"

docker rm -f bckupr
docker rm -f data_writer

docker volume rm test_volume_backup
rm -rf /tmp/bckupr/mount
