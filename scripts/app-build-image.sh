set -e

export VERSION=${VERSION:-local}
export BUILDER_IMAGE=golang:1.22-alpine
export BASE_IMAGE=scratch

docker buildx build ${DOCKER_ARGS} \
    --build-arg CREATED=$(date -u +'%Y-%m-%dT%H:%M:%S') \
    --build-arg REVISION=$(git rev-parse HEAD) \
    --build-arg BUILDER_IMAGE=${BUILDER_IMAGE} \
    --build-arg BASE_IMAGE=${BASE_IMAGE} \
    --build-arg VERSION=${VERSION} ${@:---load} \
    -t sbnarra/bckupr:${VERSION} .
echo "Built image: sbnarra/bckupr:${VERSION}"