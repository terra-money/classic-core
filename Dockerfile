FROM golang:1.13-buster AS build

RUN apt-get update
RUN apt-get install -y curl git build-essential

WORKDIR /go/src/github.com/terra-project/core

COPY . .

RUN make tools
RUN make install

RUN cp /go/pkg/mod/github.com/\!cosm\!wasm/go-cosmwasm@v*/api/libgo_cosmwasm.so /lib/libgo_cosmwasm.so


FROM ubuntu:latest

WORKDIR /root

COPY --from=build /go/bin/terrad /usr/local/bin/terrad
COPY --from=build /go/bin/terracli /usr/local/bin/terracli
COPY --from=build /lib/libgo_cosmwasm.so /lib/libgo_cosmwasm.so

CMD [ "terrad", "--help" ]
