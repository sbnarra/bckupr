set -e

export BUILDKIT_PROGRESS=plain
export VERSION=${VERSION:-local}

docker buildx build ${DOCKER_ARGS} \
    --build-arg CREATED=$(date -u +'%Y-%m-%dT%H:%M:%S') \
    --build-arg REVISION=$(git rev-parse HEAD) \
    --build-arg VERSION=${VERSION} ${@:---load} \
    -t sbnarra/bckupr:${VERSION} .
echo "Built image: sbnarra/bckupr:${VERSION}"