# Bckupr

![Build Workflow](https://github.com/sbnarra/bckupr/actions/workflows/build.yml/badge.svg)
![Dependabot Workflow](https://github.com/sbnarra/bckupr/actions/workflows/dependabot.yml/badge.svg?event=pull_request)
![Edge Workflow](https://github.com/sbnarra/bckupr/actions/workflows/edge.yml/badge.svg?branch=development)
![Latest Workflow](https://github.com/sbnarra/bckupr/actions/workflows/latest.yml/badge.svg?branch=main)

[![DockerHub Pulls](https://img.shields.io/docker/pulls/sbnarra/bckupr.svg)](https://hub.docker.com/r/sbnarra/bckupr)
![Latest Version](https://img.shields.io/docker/v/sbnarra/bckupr/latest)
![Image Size](https://img.shields.io/docker/image-size/sbnarra/bckupr/latest)

Bckupr is program to automate backup creation and data restoration.

This tool reads docker labels to determine which volumes/mounts require backing up, whilst managing containers write access and their dependant containers, creates local backups and offers options for offsite storage.

Features:
* Simple label configuration
* Schedule backups and rotations using cron expressions
* Builtin CLI tool for managing backups
* Web interface for managing backups
* Notifications - via shoutrr
* Metrics via Prometheus
* backup retention

See the [User Docs](https://sbnarra.github.io/bckupr) for detailed instructions on running Bckupr.

## Getting Started

The project is built using GoLang and is published as a Docker image. Python is also used for building user documentation with mkdocs though docker is used as a wrapper for testing and building.

### Prerequisites

This project requires the following tools:

* Go - https://go.dev/doc/install
* Docker - https://docs.docker.com/engine/install/

### Installation

Initiliase the project:
```shell
# Clone a local copy of the repository
git clone git@github.com:sbnarra/bckupr.git
# Change directory into the new clone
cd bckupr
# Initialise the dependancies
go mod init
# Run the project to see the CLI help menu
go run .
# Open UI
open http://localhost:8000
```

## Running the tests

The project includes unit tests written in Go and end-to-ends tests in Bash. (_should rewrite in Go within the `test` dir_)

Run Go unit tests:
```shell
go test -v ./...
```

Run End2End tests:
```shell
./scripts/app-test-end2end.sh
```

## Building the Project

The project is built automatically on each commit using GitHub actions. To publish a new image and documentation run the same pipeline manually run.

Building image locally:
```shell
./scripts/app-build-image.sh
docker run \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /tmp/backups:/tmp/backups \
    -e BACKUP_DIR=/tmp/backups \
    -p 8000:8000 \
    sbnarra/bckupr:local
```

Building user documentation locally:
```shell
./scripts/docs-build-site.sh
```

Viewing user documentation locally on port 8000:
```shell
./scripts/docs-serve-local.sh
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/sbnarra/bckupr/tags). 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
