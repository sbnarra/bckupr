ARG BASE_IMAGE
FROM ${BUILDER_IMAGE:-golang:1.22-alpine} AS base

ENV GO111MODULE=on \
    CGO_ENABLED=0
    # CGO_ENABLED=0 \
    # GOOS=linux \
    # GOARCH=amd64

WORKDIR /bckupr

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o bckupr .

FROM ${BASE_IMAGE:-scratch}

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/bckupr"
LABEL org.opencontainers.image.title "bckupr"
LABEL org.opencontainers.image.description "container backup/restore"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/bckupr"
LABEL org.opencontainers.image.documentation "https://sbnarra.github.io/bckupr"

ARG CREATED
LABEL org.opencontainers.image.created ${CREATED:-unset}
ARG VERSION
LABEL org.opencontainers.image.version ${VERSION:-unset}
ARG REVISION
LABEL org.opencontainers.image.revision ${REVISION:-unset}
ARG BASE_IMAGE
LABEL org.opencontainers.image.base.name ${BASE_IMAGE:-alpine}

WORKDIR /

COPY --from=base /bckupr/ui /ui/
COPY --from=base /bckupr/bckupr /

COPY configs/offsite /offsite
COPY configs/local /local
ENV LOCAL_CONTAINERS=/local/tar.yml

ENTRYPOINT ["/bckupr"]
CMD ["daemon"]

EXPOSE 8000
VOLUME /var/run/docker.sock