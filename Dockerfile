ARG ALPINE_VERSION=3.15
ARG GO_VERSION=1.17.8

# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS go-builder

# install deps
RUN set -eux; \
    apk add --no-cache ca-certificates build-base git cmake

WORKDIR /code
COPY . /code/

# install mimalloc
RUN git clone --depth 1 https://github.com/microsoft/mimalloc; \
    cd mimalloc; \
    mkdir build; \
    cd build; \
    cmake ..; \
    make -j$(nproc); \
    make install
ENV MIMALLOC_RESERVE_HUGE_OS_PAGES=4

# install libwasmvm; see https://github.com/CosmWasm/wasmvm/releases
ARG WASMVM_VERSION=0.16.6
RUN case $(uname -m) in \
        x86_64 | amd64) \
            ARCH=x86_64 \
            WASMVM_SHA256=fe63ff6bb75cad9116948d96344391d6786b6009d28e7016a85e1a268033d8f8;; \
        aarch64 | arm64) \
            ARCH=aarch64; \
            WASMVM_SHA256=dda9376d437cc8e0b9f325621887454a29660627a61c93841689338557494b50;; \
        *) echo "Unkown architecture" && exit 1;; \
    esac; \
    wget https://github.com/CosmWasm/wasmvm/releases/download/v${WASMVM_VERSION}/libwasmvm_muslc.${ARCH}.a -O /lib/libwasmvm_muslc.a; \
    echo "${WASMVM_SHA256} */lib/libwasmvm_muslc.a" | sha256sum -c -

# run build and force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LDFLAGS="-linkmode=external -extldflags \"-L/code/mimalloc/build -lmimalloc -Wl,-z,muldefs -static\"" make build

FROM alpine:${ALPINE_VERSION} AS runtime

RUN addgroup terra && \
    adduser -G terra -D -h /terra terra

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
