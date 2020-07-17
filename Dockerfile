# Simple usage with a mounted data directory:
# > docker build -t terramoney/core .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terramoney/core terrad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terramoney/core terrad start
FROM cosmwasm/go-ext-builder:0.8.2-alpine AS rust-builder

RUN apk add git

# copy go dependency files into go project path
WORKDIR /go/src/github.com/terra-project/core
COPY go.* /go/src/github.com/terra-project/core/

# downlaod go-cosmwasm
RUN go mod download github.com/CosmWasm/go-cosmwasm

RUN export GO_WASM_DIR=$(go list -f "{{ .Dir }}" -m github.com/CosmWasm/go-cosmwasm) && \
     cd ${GO_WASM_DIR} && \
     cargo build --release --features backtraces --example muslc && \
     mv ${GO_WASM_DIR}/target/release/examples/libmuslc.a /lib/libgo_cosmwasm_muslc.a

# --------------------------------------------------------
FROM cosmwasm/go-ext-builder:0.8.2-alpine AS go-builder

RUN apk add git
# without this, build with LEDGER_ENABLED=false
RUN apk add libusb-dev linux-headers

WORKDIR /go/src/github.com/terra-project/core

COPY . .

# Copy shared library from rust-builder
COPY --from=rust-builder /lib/libgo_cosmwasm_muslc.a /lib/libgo_cosmwasm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN BUILD_TAGS=muslc make build

# --------------------------------------------------------
FROM alpine:3.12

COPY --from=go-builder /go/src/github.com/terra-project/core/build/terrad /usr/bin/terrad
COPY --from=go-builder /go/src/github.com/terra-project/core/build/terracli /usr/bin/terracli

WORKDIR /root

CMD [ "terrad", "--help" ]
