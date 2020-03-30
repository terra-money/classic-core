# Simple usage with a mounted data directory:
# > docker build -t terra .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad start
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.terrad:/root/.terrad -v ~/.terracli:/root/.terracli terra terrad start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/terra

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make tools && \
    make go-mod-cache && \
    make build-linux && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates rsync jq curl

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/terrad /usr/bin/terrad
COPY --from=build-env /go/bin/terracli /usr/bin/terracli

# Create a terra group and a terra user
RUN addgroup -S terra -g 54524 && adduser -S terra -u 54524 -h /home/terra -G terra

# Tell docker that all future commands should run as the terra user
USER terra
WORKDIR /home/terra

# Run terrad by default, omit entrypoint to ease using container with terracli
CMD ["terrad"]
