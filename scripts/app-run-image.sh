set -e

docker run --rm -it --name bckupr \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /tmp/backups:/backups \
    sbnarra/bckupr:${VERSION:-local} ${@:-sh}