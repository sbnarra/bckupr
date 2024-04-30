ARG BASE_IMAGE
FROM ${BUILDER_IMAGE:-golang:1.22-alpine} AS base

RUN mkdir /empty-dir

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app
COPY ./ /app/.
RUN go build -o bckupr .

FROM ${BASE_IMAGE:-scratch}

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/bckupr"
LABEL org.opencontainers.image.title "bckupr"
LABEL org.opencontainers.image.description "docker volumes backup/restore manager"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/bckupr"
LABEL org.opencontainers.image.documentation "https://sbnarra.github.io/bckupr"

ARG CREATED
LABEL org.opencontainers.image.created ${CREATED:-unset}
ARG VERSION
LABEL org.opencontainers.image.version ${VERSION:-unset}
ARG REVISION
LABEL org.opencontainers.image.revision ${REVISION:-unset}
ARG BASE_IMAGE
LABEL org.opencontainers.image.base.name ${BASE_IMAGE:-scratch}

ENV VERSION ${VERSION:-unset}
ENV RUNNING_IN_CONTAINER 1

WORKDIR /

COPY --from=base /app/bckupr /bin/bckupr
COPY --from=base /empty-dir /tmp
ENV PATH /bin

COPY web/ /web
ENV UI_BASE_PATH /

COPY configs/local/ /local
ENV LOCAL_CONTAINERS_CONFIG=/local/tar.yml

COPY configs/offsite/ /offsite

COPY configs/rotation /rotation
ENV ROTATION_POLICIES_CONFIG=/rotation/policies.yaml

ENTRYPOINT ["bckupr"]

EXPOSE 8000
VOLUME /var/run/docker.sock