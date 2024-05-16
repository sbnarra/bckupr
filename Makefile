.DEFAULT_GOAL=build

VERSION?=local
BUILD_ARGS?=--load

ARGS?=--rm -it
CMD?=sh
BACKUP_DIR?=/tmp/backups-${VERSION}

DOCS_PATH=docs/gh-pages

clean:
	go clean ./...
	sudo rm -rf bckupr ${DOCS_PATH}/site
	docker compose -f docker-compose.dev.yml down || true

generate:
	go generate ./...
	make ui-build

test:
	go test -p 1 -v ./...
	./scripts/app-test-end2end.sh

build: test
	go build

run:
	go run . daemon --host-backup-dir ${BACKUP_DIR} --container-backup-dir ${BACKUP_DIR}

package: generate
	docker buildx build ${BUILD_ARGS} \
    	--build-arg CREATED=$(shell date -u +'%Y-%m-%dT%H:%M:%S') \
    	--build-arg REVISION=$(shell git rev-parse HEAD) \
    	--build-arg VERSION=${VERSION} \
    	-t sbnarra/bckupr:${VERSION} .

package-run: package
	mkdir -p ${BACKUP_DIR} # want to maintain user permissions so pre-creating
	docker run --name bckupr ${ARGS} \
    	-v /var/run/docker.sock:/var/run/docker.sock \
    	-v ${BACKUP_DIR}:/backups \
    	sbnarra/bckupr:${VERSION} ${CMD}

generate-docs:
	docker run --rm \
		-v ${PWD}:/bckupr:rw -w /bckupr/${DOCS_PATH} \
		python:3.9-slim \
		sh -c "pip install -r requirements.txt && mkdocs build --config-file mkdocs.yml"

run-docs:
	docker run --rm -it \
		-v ${PWD}:/bckupr:ro -w /bckupr/${DOCS_PATH} \
		-p 8000:8000 \
		python:3.9-slim \
		sh -c "pip install -r requirements.txt && mkdocs serve --config-file mkdocs.yml"

ui-build:
	docker run --rm -it \
		-v ${PWD}:/bckupr:rw -w /bckupr/web/ \
		node:20-alpine \
		sh -c "npm install && npm run build"

dev:
	docker compose -f docker-compose.dev.yml up