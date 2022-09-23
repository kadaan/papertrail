# syntax=docker/dockerfile:experimental

ARG VERSION
ARG REVISION
ARG USER
ARG HOST
ARG BUILD_DATE
ARG BRANCH

FROM golang:1.19.1-alpine AS base
ARG VERSION
ARG REVISION
ARG USER
ARG HOST
ARG BUILD_DATE
ARG BRANCH
ENV VERSION=$VERSION REVISION=$REVISION USER=$USER HOST=$HOST BUILD_DATE=$BUILD_DATE BRANCH=$BRANCH
WORKDIR /src
RUN --mount=type=cache,id=apk,sharing=locked,target=/var/cache/apk ln -vs /var/cache/apk /etc/apk/cache && \
    apk add --update git gcc libc-dev && \
    mkdir /archives && \
    mkdir /dist
COPY . .
WORKDIR /src/lib/web/ui
RUN go generate
WORKDIR /src

FROM base as darwin
RUN GOOS=darwin GARCH=amd64 go build \
            -o /dist/papertrail_darwin \
            -a \
            -ldflags "-s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/papertrail/version.Version=$VERSION -X github.com/kadaan/papertrail/version.Revision=$REVISION -X github.com/kadaan/papertrail/version.Branch=$BRANCH -X github.com/kadaan/papertrail/version.BuildUser=$USER -X github.com/kadaan/papertrail/version.BuildHost=$HOST -X github.com/kadaan/papertrail/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo'
FROM base as linux
RUN GOOS=linux GARCH=amd64 go build \
            -o /dist/papertrail_linux \
            -a \
            -ldflags "-d -s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/papertrail/version.Version=$VERSION -X github.com/kadaan/papertrail/version.Revision=$REVISION -X github.com/kadaan/papertrail/version.Branch=$BRANCH -X github.com/kadaan/papertrail/version.BuildUser=$USER -X github.com/kadaan/papertrail/version.BuildHost=$HOST -X github.com/kadaan/papertrail/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo netgo static_build' \
            -installsuffix netgo

FROM scratch as artifact
COPY --from=darwin /dist/ ./dist/
COPY --from=linux /dist/ ./dist/