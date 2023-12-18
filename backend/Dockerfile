# syntax = docker/dockerfile:1-experimental

ARG GO_VERSION=1.20
ARG PROTOBUF_VERSION=22.0
ARG APP_VERSION
ARG BUILD_DATE

#############################################################################
FROM golang:${GO_VERSION}-alpine AS base

RUN apk add build-base
ENV CGO_ENABLED=1

# Add a user to run in non-root mode
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /app

COPY go.* .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

#############################################################################
FROM base AS builder

ARG APP_VERSION
ARG BUILD_DATE

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

WORKDIR /app

COPY . .

# Build MEx Metadata service
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.build=${APP_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o mex-metadata-svc \
    ./mex/services/metadata/cmd

# Build MEx Query service
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.build=${APP_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o mex-query-svc \
    ./mex/services/query/cmd

# Build MEx Index service
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.build=${APP_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o mex-index-svc \
    ./mex/services/index/cmd

# Build MEx Auth service
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.build=${APP_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o mex-auth-svc \
    ./mex/services/auth/cmd

# Build MEx Config service
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.build=${APP_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o mex-config-svc \
    ./mex/services/config/cmd

#############################################################################
# UNIT-TEST stage: execute unit tests from the code stage results
FROM base AS unit-test
ENV CGO_ENABLED=1
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go test -v ./mex/...

#############################################################################
FROM golang:${GO_VERSION}-alpine AS final

RUN apk update && apk upgrade && apk --update add libssl3 libcrypto3

WORKDIR /app

COPY --from=builder /app/mex-metadata-svc .
COPY --from=builder /app/mex-query-svc    .
COPY --from=builder /app/mex-index-svc    .
COPY --from=builder /app/mex-auth-svc     .
COPY --from=builder /app/mex-config-svc   .


EXPOSE 3000

USER nobody:nobody

ENTRYPOINT [ "/app/mex-metadata-svc" ]
