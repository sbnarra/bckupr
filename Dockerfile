FROM --platform=$BUILDPLATFORM node:20.13.1-alpine3.19 AS node
WORKDIR /web
COPY ./web/ .
RUN npm install && npm run build

FROM --platform=$BUILDPLATFORM golang:1.22.3-alpine3.19 AS go
WORKDIR /
COPY ./ .
ARG TARGETOS
ARG TARGETARCH
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bckupr .

FROM alpine:3.20
# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/bckupr"
LABEL org.opencontainers.image.title "bckupr"
LABEL org.opencontainers.image.description "docker volumes backup/restore manager"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/bckupr"
LABEL org.opencontainers.image.documentation "https://sbnarra.github.io/bckupr"
ARG CREATED
ENV CREATED ${CREATED:-unset}
LABEL org.opencontainers.image.created ${CREATED:-unset}
ARG VERSION
ENV VERSION ${VERSION:-unset}
LABEL org.opencontainers.image.version ${VERSION:-unset}
ARG REVISION
LABEL org.opencontainers.image.revision ${REVISION:-unset}
ARG BASE_IMAGE
LABEL org.opencontainers.image.base.name ${BASE_IMAGE:-scratch}

ARG S6_OVERLAY_VERSION=3.1.6.2
# https://github.com/just-containers/s6-overlay/tree/master?tab=readme-ov-file#which-architecture-to-use-depending-on-your-targetarch
ARG TARGETARCH
RUN ARCH=$TARGETARCH && \
    ARCH=$([ $TARGETARCH == "amd64" ] && echo x86_64 || echo $ARCH) && \
    ARCH=$([ $TARGETARCH == "arm64" ] && echo aarch64 || echo $ARCH) && \
    ARCH=$([ $TARGETARCH == "386" ] && echo i686 || echo $ARCH) && \
    ARCH=$([ $TARGETARCH == "ppc64le" ] && echo powerpc64le || echo $ARCH) && \
    echo TARGETARCH=${TARGETARCH} ARCH=${ARCH} && \
    wget https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz -O /tmp/s6-overlay-noarch.tar.xz && \
        tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && \
    wget https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${ARCH}.tar.xz -O /tmp/s6-overlay-${ARCH}.tar.xz && \
        tar -C / -Jxpf /tmp/s6-overlay-${ARCH}.tar.xz && \
    rm -rf /tmp/*.tar.xz

COPY configs/s6-rc.d /etc/s6-overlay/s6-rc.d/bckupr
ENV S6_KEEP_ENV=1
RUN touch /etc/s6-overlay/s6-rc.d/user/contents.d/bckupr

ENTRYPOINT ["/init"]
RUN echo 'while [ "0" != $(ps | grep -v grep | grep "s6-supervise bckupr" | wc -l) ]; do sleep 1; done' > /cmdless
RUN chmod +x /cmdless
CMD ["/cmdless"]

COPY configs/local/ /opt/bckupr/config/local

COPY configs/offsite/ /opt/bckupr/config/offsite
RUN ln -s /opt/bckupr/config/offsite /offsite

COPY configs/rotation /opt/bckupr/config/rotation
RUN ln -s /opt/bckupr/config/rotation /rotation

ENV LOCAL_CONTAINERS_CONFIG=/opt/bckupr/config/local/tar.yml
ENV ROTATION_POLICIES_CONFIG=/opt/bckupr/config/rotation/policies.yaml

ENV UI_BUNDLE /opt/bckupr/web
ENV BCKUPR_IN_CONTAINER 1
ENV GIN_MODE release

COPY --from=node /web/out /opt/bckupr/web
COPY --from=go /bckupr /opt/bckupr/bckupr
RUN ln -s /opt/bckupr/bckupr /bin/bckupr

EXPOSE 8000
VOLUME /var/run/docker.sock
VOLUME /backups