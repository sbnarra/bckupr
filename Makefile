.DEFAULT_GOAL=build

clean:
	go clean ./...
	sudo rm -rf bckupr docs/gh-site/site

docs:
	bash ./scripts/mkdocs-build.sh
	bash ./scripts/docs-build-site.sh

generate:
	go generate ./...

test: generate
	go test -p 1 -v ./...
	./scripts/app-test-end2end.sh

build: test
	go build

VERSION?=local
PACKAGE_RESULT?=--load # --push
package: 
	docker buildx build ${DOCKER_ARGS} \
    	--build-arg CREATED=$(shell date -u +'%Y-%m-%dT%H:%M:%S') \
    	--build-arg REVISION=$(shell git rev-parse HEAD) \
    	--build-arg VERSION=${VERSION} ${PACKAGE_RESULT} \
    	-t sbnarra/bckupr:${VERSION} .
