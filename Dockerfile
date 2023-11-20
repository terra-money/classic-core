# syntax=docker/dockerfile:1

ARG source=./
ARG GO_VERSION="1.20"
ARG BUILDPLATFORM=linux/amd64
ARG BASE_IMAGE="golang:${GO_VERSION}-alpine"
FROM --platform=${BUILDPLATFORM} ${BASE_IMAGE} as base

###############################################################################
# Builder
###############################################################################

FROM base as builder-stage-1

ARG source
ARG GIT_COMMIT
ARG GIT_VERSION
ARG BUILDPLATFORM
ARG GOOS=linux \
    GOARCH=amd64

ENV GOOS=$GOOS \ 
    GOARCH=$GOARCH

# NOTE: add libusb-dev to run with LEDGER_ENABLED=true
RUN set -eux &&\
    apk update &&\
    apk add --no-cache \
    ca-certificates \
    linux-headers \
    build-base \
    cmake \
    git

# install mimalloc for musl
WORKDIR ${GOPATH}/src/mimalloc
RUN set -eux &&\
    git clone --depth 1 --branch v2.1.2 \
        https://github.com/microsoft/mimalloc . &&\
    mkdir -p build &&\
    cd build &&\
    cmake .. &&\
    make -j$(nproc) &&\
    make install

# download dependencies to cache as layer
WORKDIR ${GOPATH}/src/app
COPY ${source}go.mod ${source}go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download -x

# Cosmwasm - Download correct libwasmvm version and verify checksum
RUN set -eux &&\
    WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm | cut -d ' ' -f 5) && \
    WASMVM_DOWNLOADS="https://github.com/classic-terra/wasmvm/releases/download/${WASMVM_VERSION}"; \
    wget ${WASMVM_DOWNLOADS}/checksums.txt -O /tmp/checksums.txt; \
    if [ ${BUILDPLATFORM} = "linux/amd64" ]; then \
        WASMVM_URL="${WASMVM_DOWNLOADS}/libwasmvm_muslc.x86_64.a"; \
    elif [ ${BUILDPLATFORM} = "linux/arm64" ]; then \
        WASMVM_URL="${WASMVM_DOWNLOADS}/libwasmvm_muslc.aarch64.a"; \      
    else \
        echo "Unsupported Build Platfrom ${BUILDPLATFORM}"; \
        exit 1; \
    fi; \
    wget ${WASMVM_URL} -O /lib/libwasmvm_muslc.a; \
    CHECKSUM=`sha256sum /lib/libwasmvm_muslc.a | cut -d" " -f1`; \
    grep ${CHECKSUM} /tmp/checksums.txt; \
    rm /tmp/checksums.txt

###############################################################################

FROM builder-stage-1 as builder-stage-2

ARG source
ARG GOOS=linux \
    GOARCH=amd64

ENV GOOS=$GOOS \ 
    GOARCH=$GOARCH

# Copy the remaining files
COPY ${source} .

# Build app binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go install \
        -mod=readonly \
        -tags "netgo,muslc" \
        -ldflags " \
            -w -s -linkmode=external -extldflags \
            '-L/go/src/mimalloc/build -lmimalloc -Wl,-z,muldefs -static' \
            -X github.com/cosmos/cosmos-sdk/version.Name='terra' \
            -X github.com/cosmos/cosmos-sdk/version.AppName='terrad' \
            -X github.com/cosmos/cosmos-sdk/version.Version=${GIT_VERSION} \
            -X github.com/cosmos/cosmos-sdk/version.Commit=${GIT_COMMIT} \
            -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,muslc' \
        " \
        -trimpath \
        ./...

################################################################################

FROM alpine as terra-core

RUN apk update && apk add wget lz4 aria2 curl jq gawk coreutils "zlib>1.2.12-r2" "libssl1.1>1.1.1q-r0"

COPY --from=builder-stage-2 /go/bin/terrad /usr/local/bin/terrad

RUN addgroup -g 1000 terra && \
    adduser -u 1000 -G terra -D -h /app terra

# rest server
EXPOSE 1317
# grpc server
EXPOSE 9090
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

WORKDIR /app

CMD ["terrad", "version"]