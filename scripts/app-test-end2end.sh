TEST_DIR="$(cd $(dirname $0)/..; pwd; cd - >/dev/null)/.test_filesystem"

check_data_is() {
    _check_data_is "volume" $@
    _check_data_is "mount" $@
    echo "data matches: '$*'"
}

_check_data_is() {
    vol=$1; shift
    MSG=$(docker exec data_writer cat /mnt/${vol}/msg)
    if [ "$MSG" != "$*" ]; then
        echo "invalid mount data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
}

write_data() {
    docker exec data_writer sh -c "echo $1 >/mnt/volume/msg"
    docker exec data_writer sh -c "echo $1 >/mnt/mount/msg"
    echo "Wrote data '$1'"
}

cleanup_resources() {
    docker rm -f bckupr
    docker rm -f data_writer
    docker volume rm test_volume_backup
    rm -rf $TEST_DIR/example-mount
}

### pre test setup ###
cleanup_resources
mkdir -p $TEST_DIR/backups

set -e ### starting test ###

docker volume create test_volume_backup

VERSION=test ./scripts/app-build-image.sh
docker run --name bckupr -d -e DRY_RUN=true -e DEBUG=true \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $TEST_DIR/backups:/backups \
    sbnarra/bckupr:test

docker run --name data_writer -d \
    -l bckupr.volumes=test_volume_backup \
    -l bckupr.volumes.test_mount_backup=$TEST_DIR/example-mount \
    -v test_volume_backup:/mnt/volume \
    -v $TEST_DIR/example-mount:/mnt/mount \
    alpine sleep 120

on_exit() {
    exit_code=$?
    set +e
    if [ "$exit_code" != "0" ]; then
        docker logs bckupr
        docker logs data_writer
    fi
    cleanup_resources

    if [ "$exit_code" != "0" ]; then
        echo "Tests have failed!"
    else
        echo "Tests have passed!"
    fi
    exit $exit_code
}
trap 'on_exit' EXIT
# exit 3

BACKUP_ID="test-$(date +%Y-%m-%d_%H-%M)"
write_data pre-backup
docker exec bckupr bckupr backup --dry-run=false --backup-id=$BACKUP_ID

write_data post-backup
docker exec bckupr bckupr restore --dry-run=false --include-names data_writer --backup-id $BACKUP_ID
check_data_is pre-backup