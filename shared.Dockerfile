FROM golang:1.17-buster AS go-builder

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apt update
RUN apt install -y curl git build-essential
# debug: for live editting in the image
RUN apt install -y vim

WORKDIR /code
COPY . /code/

RUN LEDGER_ENABLED=false make build

RUN cp /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/api/libwasmvm.so /lib/libwasmvm.so

FROM ubuntu:20.04

WORKDIR /root

COPY --from=go-builder /code/build/terrad /usr/local/bin/terrad
COPY --from=go-builder /lib/libwasmvm.so /lib/libwasmvm.so

# rest server
EXPOSE 1317
# grpc
EXPOSE 9090
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

CMD ["/usr/local/bin/terrad", "version"]
