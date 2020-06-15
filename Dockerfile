# Simple usage with a mounted data directory:
# > docker build -t terra .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad start
FROM golang:1.13-buster AS build-env

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apt-get update
RUN apt-get install -y curl git build-essential

# Set working directory for the build
WORKDIR /go/src/github.com/terra-project/core

# Add source files
COPY . .

# Install tools & install core
RUN make tools
RUN make install

# Install libgo_cosmwasm.so to a shared directory where it is readable by all users
# See https://github.com/CosmWasm/wasmd/issues/43#issuecomment-608366314
# Note that CosmWasm gets turned into !cosm!wasm in the pkg/mod cache
RUN cp /go/pkg/mod/github.com/\!cosm\!wasm/go-cosmwasm@v*/api/libgo_cosmwasm.so /lib/x86_64-linux-gnu

WORKDIR /root

# Run terrad by default, omit entrypoint to ease using container with terracli
CMD ["terrad"]
