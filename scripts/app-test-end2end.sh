TEST_CONTAINER=backup_test
BACKUP_ID="test-$(date +%Y-%m-%d_%H-%M)"
VERSION=${VERSION:-local}

check_data_is() {
    MSG=$(docker exec $TEST_CONTAINER cat /mnt/volume/msg)
    if [ "$MSG" != "$*" ]; then
        echo "invalid volume data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
    MSG=$(docker exec $TEST_CONTAINER cat /mnt/mount/msg)
    if [ "$MSG" != "$*" ]; then
        echo "invalid mount data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
    # MSG=$(docker exec $TEST_CONTAINER cat /msg)
    # if [ "$MSG" != "$*" ]; then
    #     echo "invalid mount data: on-disk='$MSG' != expected='$*'"; exit 1
    # fi
    echo "data matches: '$*'"
}

write_data() {
    docker exec $TEST_CONTAINER sh -c "echo $1 >/mnt/volume/msg"
    docker exec $TEST_CONTAINER sh -c "echo $1 >/mnt/mount/msg"
    # docker exec $TEST_CONTAINER sh -c "echo $1 >/msg"
    echo "Wrote data '$1'"
}

bckupr() {
    go run . $@ --backup-dir $PWD/.test_filesystem/backups --dry-run=false --no-daemon=true
}

# setup
rm -rf $PWD/.test_filesystem
mkdir -p $PWD/.test_filesystem/backups
rm -rf /tmp/bckupr/mount
docker rm -f $TEST_CONTAINER
docker volume rm test_volume_backup

set -e # start test

mkdir -p /tmp/bckupr/mount
docker volume create test_volume_backup

TEST_FS="$PWD/.test_filesystem"
docker run --name $TEST_CONTAINER -d \
    -l bckupr.stop=true \
    \
    -l bckupr.filesystem=false \
    \
    -l bckupr.volumes=test_volume_backup \
    -v test_volume_backup:/mnt/volume \
    \
    -l bckupr.volumes.test_mount_backup=$TEST_FS/example-mount \
    -v $TEST_FS/example-mount:/mnt/mount \
    \
    alpine sleep 120

# write intial data
write_data pre-backup
check_data_is pre-backup

# perform backup of intial data
bckupr backup --backup-id=$BACKUP_ID

# update container data
write_data post-backup
check_data_is post-backup

# restore to initial data
bckupr restore --include-names $TEST_CONTAINER --backup-id $BACKUP_ID
check_data_is pre-backup

echo "test completed successfully, running clean up"
docker rm -f $TEST_CONTAINER
docker volume rm test_volume_backup
rm -rf /tmp/bckupr/mount
