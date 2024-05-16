#!/bin/bash

check_data_is() {
    _check_data_is "volume" $@
    _check_data_is "mount" $@
    echo "data matches: '$*'"
}

_check_data_is() {
    MSG=$(docker exec "rw-$BACKUP_ID" cat /mnt/${1}/msg); shift
    if [ "$MSG" != "$*" ]; then
        echo "invalid mount data: on-disk='$MSG' != expected='$*'"
        exit 1
    fi
}

write_data() {
    docker exec "rw-$BACKUP_ID" sh -c "echo $1 | tee /mnt/volume/msg"
    docker exec "rw-$BACKUP_ID" sh -c "echo $1 | tee /mnt/mount/msg"
}

cleanup_resources() {
    docker rm -f bckupr
    docker rm -f "rw-$BACKUP_ID"
    docker volume rm test_volume_backup
    rm -rf $TEST_DIR/example-mount
}

### pre test setup ###
BACKUP_ID="$(date +%Y%m%d%H%M)-cli"
TEST_DIR="$(cd $(dirname $0)/..; pwd; cd - >/dev/null)/.test_filesystem"
cleanup_resources
mkdir -p $TEST_DIR/backups

set -e ### starting test ###
docker volume create test_volume_backup
make package-run CMD="" ARGS="-d -e DEBUG=$DEBUG" BACKUP_DIR=$TEST_DIR/backups VERSION=test
docker run --name "rw-$BACKUP_ID" -d \
    -l bckupr.volumes=test_volume_backup \
    -l bckupr.volumes.test_mount_backup=$TEST_DIR/example-mount \
    -v test_volume_backup:/mnt/volume \
    -v $TEST_DIR/example-mount:/mnt/mount \
    alpine sleep 120

on_exit() {
    exit_code=$?
    set +e
    if [ "$exit_code" != "0" ] || [ "$DEBUG" == "1" ]; then
        echo bckupr logs && docker logs bckupr
        echo "rw-$BACKUP_ID" logs && docker logs "rw-$BACKUP_ID"
    fi
    cleanup_resources
    echo "Tests have $([ "$exit_code" != "0" ] && echo failed || echo passed)!"
    exit $exit_code
}
trap 'on_exit' EXIT

write_data "pre-backup: $BACKUP_ID"
docker exec bckupr bckupr backup --no-dry-run --backup-id=$BACKUP_ID

write_data "post-backup: $BACKUP_ID"
docker exec bckupr bckupr restore --no-dry-run --include-names "rw-$BACKUP_ID" --backup-id $BACKUP_ID
check_data_is "pre-backup: $BACKUP_ID"