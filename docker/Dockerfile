ARG version=v0.5.9-oracle

FROM terramoney/core:${version}

ARG chainid=columbus-5

ENV CHAINID ${chainid}

# Moniker will be updated by entrypoint.
RUN terrad init --chain-id $chainid moniker

# Backup for templating
RUN mv ~/.terra/config/config.toml ~/config.toml
RUN mv ~/.terra/config/app.toml ~/app.toml

RUN if [ "$chainid" = "columbus-5" ] ; then wget -O ~/.terra/config/genesis.json https://columbus-genesis.s3.ap-northeast-1.amazonaws.com/columbus-5-genesis.json; fi
RUN if [ "$chainid" = "columbus-5" ] ; then wget -O ~/.terra/config/addrbook.json https://network.terra.dev/addrbook.json; fi

RUN if [ "$chainid" = "bombay-12" ] ; then wget -O ~/.terra/config/genesis.json https://raw.githubusercontent.com/terra-money/testnet/master/bombay-12/genesis.json; fi
RUN if [ "$chainid" = "bombay-12" ] ; then wget -O ~/.terra/config/addrbook.json https://raw.githubusercontent.com/terra-money/testnet/master/bombay-12/addrbook.json; fi

RUN apk update && apk add wget lz4 aria2 curl jq gawk coreutils

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

CMD ["/usr/local/bin/terrad", "start"]