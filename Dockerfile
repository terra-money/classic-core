# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:1.17-alpine3.14 AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# See https://github.com/terra-money/wasmvm/releases
ADD https://github.com/terra-money/wasmvm/releases/download/v0.16.1/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 7a44ca40fbd1e4b27325fa5053c4699b2b7a5dcb043a8763792acfa8f2adae70

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make build

FROM alpine:3.12

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
