# Simple usage with a mounted data directory:
# > docker build -t terra .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad start
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/terra

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make get_tools && \
    make get_vendor_deps && \
    make build-linux && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates rsync jq
WORKDIR /etc/terrad

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/terrad /usr/bin/terrad
COPY --from=build-env /go/bin/terracli /usr/bin/terracli

# Run terrad by default, omit entrypoint to ease using container with terracli
EXPOSE 26656 26657
CMD ["terrad"]
