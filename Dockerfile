ARG RUST_VERSION=1.55.0
ARG GO_VERSION=1.17
ARG ALPINE_VERSION=3.14
ARG BUILD_MODE=src

FROM rust:${RUST_VERSION}-alpine${ALPINE_VERSION} AS rust-builder-src
WORKDIR /code

# install deps
RUN apk add --no-cache ca-certificates musl-dev

# build libwasmvm
ARG WASMVM_VERSION=0.16.5
ARG WASMVM_SHA256=9db6995f536ecf17aad078b05cebc27a7a950ef0d36f3cb7aa54a02cb4a25833
RUN set -eux; \
    wget https://github.com/CosmWasm/wasmvm/archive/refs/tags/v${WASMVM_VERSION}.tar.gz -O wasmvm.tar.gz; \
    echo "${WASMVM_SHA256} *wasmvm.tar.gz" | sha256sum -c -; \
    tar xf wasmvm.tar.gz --strip-components=1; \
    rm wasmvm.tar.gz; \
    cd libwasmvm; \
    cargo build --release --example muslc

FROM scratch AS rust-builder-bin
WORKDIR /code/libwasmvm/target/release/examples/

FROM rust-builder-${BUILD_MODE} AS rust-builder

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS go-builder
WORKDIR /code

# install deps
RUN apk add --no-cache ca-certificates build-base git

# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

# build terrad, force static lib
COPY . /code
COPY --from=rust-builder /code/libwasmvm/target/release/examples/libmuslc.a /lib/libwasmvm_muslc.a
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make build

FROM alpine:${ALPINE_VERSION} AS release
WORKDIR /root

COPY --from=go-builder /code/build/terrad /usr/local/bin/terrad

# rest server
EXPOSE 1317
# grpc
EXPOSE 9090
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

CMD ["/usr/local/bin/terrad", "version"]
