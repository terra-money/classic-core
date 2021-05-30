# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:1.15-alpine3.12 AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# See https://github.com/terra-money/go-cosmwasm/releases
ADD https://github.com/terra-money/go-cosmwasm/releases/download/v0.10.4/libgo_cosmwasm_muslc.a /lib/libgo_cosmwasm_muslc.a
RUN sha256sum /lib/libgo_cosmwasm_muslc.a | grep 2aa7b034b9340fecaa928adf3e8c093893fd6a3986a569ce7cae7528845a0951

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make update-swagger-docs build

FROM alpine:3.12

WORKDIR /root

COPY --from=go-builder /code/build/terrad /usr/local/bin/terrad
COPY --from=go-builder /code/build/terracli /usr/local/bin/terracli

CMD ["/usr/local/bin/terrad", "version"]
