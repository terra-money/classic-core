# Changed from Alpine to Ubuntu because the keyring PR is linking to libc
# Alpine uses muslc instead of libc

FROM ubuntu:18.04

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install curl jq file

VOLUME ["/terrad"]
WORKDIR /terrad
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh
