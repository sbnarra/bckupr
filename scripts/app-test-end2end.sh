BACKUP_ID="test-$(date +%Y-%m-%d_%H-%M)"
VERSION=${VERSION:-local}

check_data_is() {
    MSG=$(docker exec test_service cat /mnt/volume/msg)
    if [ "$MSG" != "$*" ]; then
        echo "invalid volume data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
    MSG=$(docker exec test_service cat /mnt/mount/msg)
    if [ "$MSG" != "$*" ]; then
        echo "invalid mount data: on-disk='$MSG' != expected='$*'"; exit 1
    fi
    echo "data matches: '$*'"
}

write_data() {
    docker exec test_service sh -c "echo $1 >/mnt/volume/msg"
    docker exec test_service sh -c "echo $1 >/mnt/mount/msg"
    echo "Wrote data '$1'"
}

bckupr() {
    docker exec bckupr_instance /bckupr $@ --backup-dir $PWD/.test_filesystem/backups
}

# setup
mkdir -p $PWD/.test_filesystem/backups
rm -rf /tmp/bckupr/mount
docker rm -f bckupr_instance
docker rm -f test_service
docker volume rm test_volume_backup

set -e # start test

mkdir -p /tmp/bckupr/mount
docker volume create test_volume_backup

VERSION=test ./scripts/app-build-image.sh
docker run --name bckupr_instance -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $PWD/.test_filesystem/backups:$PWD/.test_filesystem/backups \
    sbnarra/bckupr:test --backup-dir $PWD/.test_filesystem/backups --dry-run=false

TEST_FS="$PWD/.test_filesystem"
docker run --name test_service -d \
    -l bckupr.stop=true \
    \
    -l bckupr.volumes=test_volume_backup \
    -v test_volume_backup:/mnt/volume \
    \
    -l bckupr.volumes.test_mount_backup=$TEST_FS/example-mount \
    -v $TEST_FS/example-mount:/mnt/mount \
    alpine sleep 120

write_data pre-backup
bckupr backup --backup-id=$BACKUP_ID
write_data post-backup
bckupr restore --include-names test_service --backup-id $BACKUP_ID
check_data_is pre-backup

echo "test completed successfully, running clean up"
docker rm -f bckupr_instance
docker rm -f test_service
docker volume rm test_volume_backup
rm -rf /tmp/bckupr/mount
