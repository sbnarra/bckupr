.DEFAULT_GOAL=build

VERSION?=local
BUILD_ARGS?=--load

ARGS?=--rm -it
CMD?=sh
BACKUP_DIR?=/tmp/backups-${VERSION}

DOCS_PATH=docs/gh-site

clean:
	go clean ./...
	sudo rm -rf bckupr ${DOCS_PATH}/site

generate:
	go generate ./...
test: generate
	go test -p 1 -v ./...
	./scripts/app-test-end2end.sh
build: test
	go build

package: generate
	docker buildx build ${BUILD_ARGS} \
    	--build-arg CREATED=$(shell date -u +'%Y-%m-%dT%H:%M:%S') \
    	--build-arg REVISION=$(shell git rev-parse HEAD) \
    	--build-arg VERSION=${VERSION} \
    	-t sbnarra/bckupr:${VERSION} .
run:
	docker run --name bckupr ${ARGS} \
    	-v /var/run/docker.sock:/var/run/docker.sock \
    	-v ${BACKUP_DIR}:/backups \
    	sbnarra/bckupr:${VERSION} ${CMD}

docs:
	docker run --rm \
		-v ./:/bckupr -w /bckupr/${DOCS_PATH} \
		python:3.9-slim \
		sh -c "pip install -r requirements.txt && mkdocs build --config-file mkdocs.yml"
docs-run:
	docker run --rm \
		-v ./:/bckupr -w /bckupr/${DOCS_PATH} \
		-p 8000:8000 \
		python:3.9-slim \
		sh -c "pip install -r requirements.txt && mkdocs serve --config-file mkdocs.yml"