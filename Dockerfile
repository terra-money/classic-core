# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:1.17.8-alpine3.15 AS go-builder

# See https://github.com/CosmWasm/wasmvm/releases
ENV LIBWASMVM_VERSION=0.16.3
ENV LIBWASMVM_SHA256=3fc6d5a239f3e97ac96c1a2df3006e4107ca461da4ca318bc71cfdc3e3593125

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v${LIBWASMVM_VERSION}/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep ${LIBWASMVM_SHA256}

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make build

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
