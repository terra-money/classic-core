FROM cosmwasm/go-ext-builder:0001-alpine AS rust-builder

WORKDIR /go/src/github.com/terra-project/core

COPY go.* /go/src/github.com/terra-project/core/

RUN apk add --no-cache git \
    && go mod download github.com/CosmWasm/go-cosmwasm \
    && export GO_WASM_DIR=$(go list -f "{{ .Dir }}" -m github.com/CosmWasm/go-cosmwasm) \
    && cd ${GO_WASM_DIR} \
    && cargo build --release --features backtraces --example muslc \
    && mv ${GO_WASM_DIR}/target/release/examples/libmuslc.a /lib/libgo_cosmwasm_muslc.a


FROM cosmwasm/go-ext-builder:0001-alpine AS go-builder

WORKDIR /go/src/github.com/terra-project/core

RUN apk add --no-cache git libusb-dev linux-headers

COPY . .
COPY --from=rust-builder /lib/libgo_cosmwasm_muslc.a /lib/libgo_cosmwasm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN BUILD_TAGS=muslc make update-swagger-docs build


FROM alpine:3

WORKDIR /root

COPY --from=go-builder /go/src/github.com/terra-project/core/build/terrad /usr/local/bin/terrad
COPY --from=go-builder /go/src/github.com/terra-project/core/build/terracli /usr/local/bin/terracli

CMD [ "terrad", "--help" ]
