FROM alpine:3.7
MAINTAINER YunSuk Yeo <yun@terra.money>

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME ["/terrad"]
WORKDIR /terrad
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh

