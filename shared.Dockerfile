ARG RUST_VERSION=1.60.0
ARG GO_VERSION=1.17.8
ARG BUILD_MODE=src

FROM rust:${RUST_VERSION}-buster AS rust-builder-src
WORKDIR /code

# build libwasmvm
ARG WASMVM_VERSION=0.16.6
ARG WASMVM_SHA256=3df181d7ab80e44d55ad31ebf7ee983831abe43764781c3e07b94d20b703d079
RUN set -eux; \
    wget https://github.com/CosmWasm/wasmvm/archive/refs/tags/v${WASMVM_VERSION}.tar.gz -O wasmvm.tar.gz; \
    echo "${WASMVM_SHA256} *wasmvm.tar.gz" | sha256sum -c -; \
    tar xf wasmvm.tar.gz --strip-components=1; \
    rm wasmvm.tar.gz; \
    cd libwasmvm; \
    cargo build --release

# support pre-built libwasmvm using the arg BUILD_MODE=bin
FROM scratch AS rust-builder-bin
WORKDIR /code/libwasmvm/target/release/deps
COPY libwasmvm.so ./

FROM rust-builder-${BUILD_MODE} AS rust-builder

FROM golang:${GO_VERSION}-buster AS go-builder
WORKDIR /code

# install deps
RUN apt update && \
    apt install -y curl git build-essential vim

# build terrad
COPY . /code
COPY --from=rust-builder /code/libwasmvm/target/release/deps/libwasmvm.so /lib/libwasmvm.so
RUN LEDGER_ENABLED=false make build

FROM ubuntu:20.04 AS runtime
WORKDIR /root

COPY --from=rust-builder /code/libwasmvm/target/release/deps/libwasmvm.so /lib/libwasmvm.so
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
