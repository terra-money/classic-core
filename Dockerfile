# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:1.17.8-alpine3.15 AS go-builder

# See https://github.com/terra-money/wasmvm/releases
ENV LIBWASMVM_VERSION=1.0.0-beta10
ENV LIBWASMVM_SHA256=d1be6826066e9d292cefc71ba7ca8107a7c7fbf5a241b3d7a5c5ee4fa60cb799

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git cmake
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# Install mimalloc
RUN git clone --depth 1 https://github.com/microsoft/mimalloc; cd mimalloc; mkdir build; cd build; cmake ..; make -j$(nproc); make install
ENV MIMALLOC_RESERVE_HUGE_OS_PAGES=4

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/terra-money/wasmvm/releases/download/v${LIBWASMVM_VERSION}/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep ${LIBWASMVM_SHA256}

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LDFLAGS="-linkmode=external -extldflags \"-L/code/mimalloc/build -lmimalloc -Wl,-z,muldefs -static\"" make build

FROM alpine:3.15.4

RUN addgroup terra \
    && adduser -G terra -D -h /terra terra

WORKDIR /terra

COPY --from=go-builder /code/build/terrad /usr/local/bin/terrad

USER terra

# rest server
EXPOSE 1317
# grpc
EXPOSE 9090
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

CMD ["/usr/local/bin/terrad", "version"]
